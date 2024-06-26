/*
A game for Ebitengine game jam 2024

# Copyright (C) 2024 Loïg Jezequel

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"image"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/loig/ebitenginegamejam2024/assets"
)

type tetrisLine = [gPlayAreaWidthInBlocks]int
type tetrisGrid = [gPlayAreaHeightInBlocks + gInvisibleLines]tetrisLine

// Structure for one tetris game
type tetris struct {
	area                  tetrisGrid
	currentBlock          tetrisBlock
	nextBlock             tetrisBlock
	heldBlock             tetrisBlock
	autoDownFrame         int
	autoDownFrameLimit    int
	manualDownFrame       int
	manualDownFrameLimit  int
	lrMoveFrame           int
	lrMoveFrameLimit      int
	lrFirstMoveFrame      int
	lrFirstMoveFrameLimit int
	manualMoveAllowed     bool
	numLines              int
	dropLenght            int
	deathLines            int
	// animation and lines removal handling
	toCheck                          [2]int
	toRemove                         [4]bool
	toRemoveNum                      int
	firstAvailable                   int
	removeLineAnimationFrame         int
	removeLineAnimationStep          int
	removeLineAnimationStepNumFrames int
	inAnimation                      bool
	// invisible blocks handling
	invisibleLevel int
	invisibleStep  int
	invisibleFrame int
	// count score
	score int
	// improvements
	betterRotation      bool
	canHold             bool
	life                int
	currentLife         int
	dead                bool
	deathAnimationFrame int
}

func (t *tetris) init(level int, balance balancing, speedLevel int, score int, betterRotation, canHold bool, life, currentLife int) {
	if level == 0 {
		t.area = tetrisGrid{}
		t.currentBlock = getNewBlock(tetrisBlock{id: -1}, tetrisBlock{id: -1})
		t.currentBlock.setInitialPosition()
		t.nextBlock = getNewBlock(tetrisBlock{id: -1}, tetrisBlock{id: -1})
		t.heldBlock = tetrisBlock{id: -1}
	}
	t.autoDownFrame = 0
	t.autoDownFrameLimit = gSpeeds[balance.getSpeedLevel(speedLevel)]
	t.manualDownFrame = 0
	t.manualDownFrameLimit = 4
	t.lrMoveFrame = 0
	t.lrMoveFrameLimit = 6
	t.lrFirstMoveFrame = 0
	t.lrFirstMoveFrameLimit = 15
	t.manualMoveAllowed = true
	t.numLines = 0
	t.dropLenght = 0
	t.deathLines = balance.getDeathLines()
	t.toCheck = [2]int{}
	t.toRemove = [4]bool{}
	t.toRemoveNum = 0
	t.removeLineAnimationFrame = 0
	t.removeLineAnimationStep = 0
	t.removeLineAnimationStepNumFrames = 8
	t.invisibleFrame = 0
	t.invisibleStep = maxLevelInvisibleBlocks
	t.invisibleLevel = balance.getInvisibleBlocks()
	t.score = score

	t.betterRotation = betterRotation
	t.canHold = canHold
	t.life = life
	t.currentLife = currentLife
	t.dead = false
	t.deathAnimationFrame = 0

	t.inAnimation = false
}

func (t *tetris) setUpNext() {
	t.lost()

	if t.dead {
		t.inAnimation = true
		return
	}

	futureBlock := getNewBlock(t.currentBlock, t.nextBlock)
	t.currentBlock = t.nextBlock
	t.currentBlock.setInitialPosition()
	t.nextBlock = futureBlock

	t.manualMoveAllowed = false

	t.invisibleFrame = 0
	t.invisibleStep = maxLevelInvisibleBlocks
}

func (t *tetris) update(moveDownRequest, moveLeftRequest, moveRightRequest, holdRequest, rotateLeft, rotateRight bool, level int) (playSounds [assets.NumSounds]bool) {

	if t.dead {
		playSounds[assets.SoundDeathID] = t.deathAnimationFrame == 0
		t.deathAnimationFrame++
		if t.deathAnimationFrame >= 90 {
			t.inAnimation = false
		}
		return
	}

	if t.removeLineAnimationStep > 0 {

		t.removeLineAnimationFrame++
		if t.removeLineAnimationFrame >= t.removeLineAnimationStepNumFrames {
			t.removeLineAnimationStep++
			t.removeLineAnimationFrame = 0
		}

		if t.removeLineAnimationStep == 4 && t.removeLineAnimationFrame <= 0 {
			switch t.toRemoveNum {
			case 1:
				t.score += 40 * (level + 1)
			case 2:
				t.score += 100 * (level + 1)
			case 3:
				t.score += 300 * (level + 1)
			case 4:
				t.score += 1200 * (level + 1)
			}
			t.numLines += t.toRemoveNum
		}

		if t.removeLineAnimationStep < 8 {
			return
		}

		t.removeLineAnimationStep = 0

		// lines removal animation and effects
		playSounds[assets.SoundLinesFallingID] = true
		t.removeLines()

		t.toRemove = [4]bool{}
		t.toRemoveNum = 0
		t.toCheck = [2]int{}
		t.inAnimation = false

		t.setUpNext()

		return
	}

	if t.canHold && holdRequest {
		if canReplace(t.currentBlock.x, t.currentBlock.y, t.heldBlock, t.nextBlock, t.area) {
			t.heldBlock, t.currentBlock = t.currentBlock, t.heldBlock
			if t.currentBlock.id < 0 {
				t.currentBlock = t.nextBlock
				t.nextBlock = getNewBlock(t.heldBlock, t.currentBlock)
			}
			t.currentBlock.x = t.heldBlock.x
			t.currentBlock.y = t.heldBlock.y
			t.heldBlock.x = 0
			t.heldBlock.y = 0
		}
	}

	t.invisibleFrame++
	if t.invisibleFrame >= gInvisibleNumFrames {
		t.invisibleStep--
		t.invisibleFrame = 0
		if t.invisibleStep <= 0 {
			t.invisibleStep = maxLevelInvisibleBlocks
		}
	}

	effectiveRotation := false

	if rotateLeft && !rotateRight {
		effectiveRotation = t.currentBlock.rotateLeft(t.area)
	}

	if rotateRight && !rotateLeft {
		effectiveRotation = t.currentBlock.rotateRight(t.area)
	}

	playSounds[assets.SoundRotationID] = effectiveRotation
	if effectiveRotation && t.betterRotation {
		t.autoDownFrame = 0
	}

	mayAllowManualMoves := false

	// left/right movements of blocks handling
	xMove := 0
	if moveLeftRequest {
		xMove--
	}
	if moveRightRequest {
		xMove++
	}

	if !moveLeftRequest && !moveRightRequest {
		mayAllowManualMoves = true
		t.lrMoveFrame = 0
		t.lrFirstMoveFrame = 0
	}

	if !t.manualMoveAllowed {
		xMove = 0
	}

	if xMove != 0 {
		if t.lrMoveFrame > 0 || (t.lrFirstMoveFrame > 0 && t.lrFirstMoveFrame < t.lrFirstMoveFrameLimit) {
			xMove = 0
		}
		t.lrMoveFrame++
		if t.lrMoveFrame >= t.lrMoveFrameLimit {
			t.lrMoveFrame = 0
		}
		if t.lrFirstMoveFrame < t.lrFirstMoveFrameLimit {
			t.lrFirstMoveFrame++
		}
	}

	// automatic down movement of blocks handling
	autoDown := false
	t.autoDownFrame++

	if t.autoDownFrame >= t.autoDownFrameLimit {
		autoDown = true
		t.autoDownFrame = 0
	}

	// manual down movement of blocks handling
	manualDown := false

	if !moveDownRequest {
		t.manualDownFrame = 0
		t.manualMoveAllowed = t.manualMoveAllowed || mayAllowManualMoves
		t.dropLenght = 0
	}

	if moveDownRequest && t.manualMoveAllowed {
		manualDown = t.manualDownFrame == 0
		t.manualDownFrame++
		if t.manualDownFrame >= t.manualDownFrameLimit {
			t.manualDownFrame = 0
		}
	}

	if manualDown {
		t.dropLenght++
	}

	// update position according to movements requests
	var stuck bool
	stuck, playSounds[assets.SoundLeftRightID] = t.currentBlock.updatePosition(xMove, autoDown || manualDown, t.area)
	if stuck {
		playSounds[assets.SoundTouchGroundID] = true

		t.toCheck = t.currentBlock.writeInGrid(&t.area)

		t.score += t.dropLenght

		t.toRemoveNum, t.firstAvailable, t.toRemove = t.checkLines()

		if t.toRemoveNum > 0 {
			t.removeLineAnimationStep = 1
			t.inAnimation = true
			playSounds[assets.SoundLinesVanishingID] = true
			return
		}

		t.setUpNext()
	}

	return
}

// check if the lines in toCheck are complete
// if so, remove them and update the grid
func (t tetris) checkLines() (toRemoveNum int, firstAvailable int, toRemove [4]bool) {

	count := -1
	firstAvailable = t.toCheck[0] - 1

	// get the lines that will disapear
CheckLoop:
	for l := t.toCheck[0]; l <= t.toCheck[1]; l++ {
		count++
		for x := 0; x < len(t.area[l]); x++ {
			if t.area[l][x] == 0 {
				firstAvailable = l
				continue CheckLoop
			}
		}
		toRemove[count] = true
		toRemoveNum++
	}

	return
}

func (t *tetris) removeLines() {

	// remove them from the grid from bottom to top

	// in the removal zone
	for y := t.toCheck[1]; y >= t.toCheck[0]; y-- {
		if t.firstAvailable >= 0 {
			t.area[y] = t.area[t.firstAvailable]
			t.firstAvailable--
			for t.firstAvailable >= t.toCheck[0] && t.toRemove[t.firstAvailable-t.toCheck[0]] {
				t.firstAvailable--
			}
		} else {
			t.area[y] = tetrisLine{}
		}
	}

	// above the removal zone
	for y := t.toCheck[0] - 1; y >= 0; y-- {
		if t.firstAvailable >= 0 {
			t.area[y] = t.area[t.firstAvailable]
			t.firstAvailable--
		} else {
			t.area[y] = tetrisLine{}
		}
	}

}

// check if there is anything in the above area
// which would mean that the game is lost
func (t *tetris) lost() {
	t.currentLife = t.life
	for _, line := range t.area[:gInvisibleLines+t.deathLines] {
		for _, v := range line {
			if v != 0 {
				t.currentLife--
				if t.currentLife < 0 {
					t.dead = true
					return
				}
			}
		}
	}
}

func getNewBlock(current, next tetrisBlock) (block tetrisBlock) {

	getRandomBlock := func() tetrisBlock {
		switch rand.Intn(7) {
		case 0:
			return getIBlock()
		case 1:
			return getOBlock()
		case 2:
			return getJBlock()
		case 3:
			return getLBlock()
		case 4:
			return getSBlock()
		case 5:
			return getTBlock()
		default:
			return getZBlock()
		}
	}

	if current.id < 0 || next.id < 0 {
		return getRandomBlock()
	}

	count := 0
	block = getRandomBlock()

	for count < 2 {

		if current.id|next.id|block.id != next.id {
			return block
		}

		block = getRandomBlock()
		count++
	}

	return block

}

func (t tetris) drawHold(screen *ebiten.Image, gray uint8) {

	x := gWidth - 3*gHoldSide/4 - gPlayAreaSide
	y := gHeight - gNextBoxSide - gHoldSide/2 + 10

	options := ebiten.DrawImageOptions{}
	options.ColorScale.ScaleWithColor(color.Gray{gray})
	options.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(assets.ImageHold, &options)

	t.heldBlock.draw(screen, gray, x+gSquareSideSize/2, y+gSquareSideSize/2, 0.5)

}

func (t tetris) drawLife(screen *ebiten.Image, gray uint8) {

	if t.life > 0 {
		x := gPlayAreaSide + gPlayAreaWidth + gPlayAreaSide + gInfoLeftSide + (gInfoWidth-gHeartWidth*t.life)/2
		y := gYLevelFromTop - 2*gHeartWidth - 20

		options := ebiten.DrawImageOptions{}
		options.ColorScale.ScaleWithColor(color.Gray{gray})
		options.GeoM.Translate(float64(x), float64(y))
		for i := 0; i < t.life; i++ {
			image := assets.ImageHeart
			if i < t.currentLife {
				image = assets.ImageFullHeart
			}
			screen.DrawImage(image, &options)
			options.GeoM.Translate(float64(gHeartWidth), 0)
		}
	}

}

func (t tetris) draw(screen *ebiten.Image, gray uint8) {

	t.drawLife(screen, gray)

	xNextOrigin := gPlayAreaSide + gPlayAreaWidth + gPlayAreaSide + gInfoLeftSide + gNextMargin
	yNextOrigin := gInfoTop + gInfoSmallBoxHeight + gScoreToLevel + gInfoBoxHeight + gLevelToLines + gInfoBoxHeight + gLinesToNext + gNextMargin

	t.nextBlock.draw(screen, gray, xNextOrigin, yNextOrigin, 1)

	if t.canHold {
		t.drawHold(screen, gray)
	}

	xOrigin := gPlayAreaSide
	yOrigin := gSquareSideSize * -gInvisibleLines

	if t.removeLineAnimationStep == 0 {
		if t.invisibleStep > t.invisibleLevel || t.currentBlock.y < gInvisibleLines {
			t.currentBlock.draw(screen, gray, xOrigin, yOrigin, 1)
		}
	}

	for y, line := range t.area {
		for x, style := range line {
			if style != noStyle {

				// removal animation
				if t.removeLineAnimationStep%2 == 1 {
					if y >= t.toCheck[0] && y <= t.toCheck[1] &&
						t.toRemove[y-t.toCheck[0]] {
						if t.removeLineAnimationStep == 7 {
							continue
						}
						style = breakStyle
					}
				}

				options := ebiten.DrawImageOptions{}
				options.ColorScale.ScaleWithColor(color.Gray{gray})
				options.GeoM.Translate(float64(xOrigin+x*gSquareSideSize), float64(yOrigin+y*gSquareSideSize))
				screen.DrawImage(assets.ImageSquares.SubImage(image.Rect((style-1)*gSquareSideSize, 0, style*gSquareSideSize, gSquareSideSize)).(*ebiten.Image), &options)
			}
		}
	}

}
