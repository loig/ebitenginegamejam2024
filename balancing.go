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
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/loig/ebitenginegamejam2024/assets"
)

const (
	balanceGoalLines int = iota
	balanceSpeed
	balanceHiddenLines
	balanceDeathLines
	balanceInvisibleBlocks
	numBalances
)

const (
	maxLevelGoalLines       = 2
	maxLevelSpeed           = 5
	maxLevelHiddenLines     = 5
	maxLevelDeathLines      = 5
	maxLevelInvisibleBlocks = 3
)

type balancing struct {
	levels          [numBalances]int
	maxLevels       [numBalances]int
	choice          int
	choiceDirection int
	choices         []int
	numChoices      int
	inTransition    bool
	transitionFrame int
}

func (b *balancing) update() (end bool, playSounds [assets.NumSounds]bool) {

	if b.inTransition {
		b.transitionFrame++
		if b.transitionFrame >= gChoiceSelectionNumFrame {
			b.inTransition = false
			b.transitionFrame = 0
			if b.choiceDirection < 0 {
				b.choice = (b.choice + 1) % b.numChoices
			} else {
				b.choice = (b.choice + b.numChoices - 1) % b.numChoices
			}
		}
		return
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		playSounds[assets.SoundMenuMoveID] = true
		b.choiceDirection = 1
		b.inTransition = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		playSounds[assets.SoundMenuMoveID] = true
		b.choiceDirection = -1
		b.inTransition = true
	}

	end = inpututil.IsKeyJustPressed(ebiten.KeyEnter)

	if end {
		b.setChoice(b.choices[b.choice])
		b.choice = 0
		playSounds[assets.SoundMenuConfirmID] = true
	}

	return
}

func drawLevel(screen *ebiten.Image, level, levelMax int, x, y float64) {

	factor := 0.7
	size := factor * float64(gChoiceLevelSize)

	space := 0.1
	dx := float64(size) * (space + 1)

	steps := levelMax - levelMax/2
	if levelMax <= 3 {
		steps = levelMax
	}
	centerShift := (float64(gChoiceSize) - float64(steps)*dx - space) / 2
	yShift := float64(gChoiceSize) - 2.2*dx

	options := ebiten.DrawImageOptions{}
	options.GeoM.Scale(factor, factor)
	options.GeoM.Translate(x-dx+centerShift, y+yShift)

	for i := 0; i < levelMax; i++ {
		if i == levelMax-levelMax/2 && levelMax > 3 {
			// secondLine
			shiftAdjust := -float64(size) * 0.6
			if levelMax%2 == 0 {
				shiftAdjust = 0
			}

			options.GeoM.Translate(-dx*float64(levelMax/2)+shiftAdjust, float64(size)*0.8)
		}
		options.GeoM.Translate(dx, 0)
		shift := 0
		if i > level {
			shift = 1
		}
		screen.DrawImage(assets.ImageLevel.SubImage(image.Rect(shift*gChoiceLevelSize, 0, (shift+1)*gChoiceLevelSize, gChoiceLevelSize)).(*ebiten.Image), &options)
	}
}

func (b balancing) drawChoices(screen *ebiten.Image, cX, cY int) {

	r := float64(gHeight / 7)
	var gray uint8 = 200
	currentGray := gray

	angleShift := float64(b.choiceDirection) * float64(b.transitionFrame) / float64(gChoiceSelectionNumFrame) * (math.Pi * 2) / float64(b.numChoices)

	if !b.inTransition {
		angleShift = 0
		currentGray = 255
	}

	currentX, currentY := math.Cos(math.Pi/2+angleShift)*r, -math.Sin(math.Pi/2+angleShift)*r
	currentX += float64(cX - gChoiceSize/2)
	currentY += float64(cY - gChoiceSize/2)

	// current choice
	options := ebiten.DrawImageOptions{}
	options.ColorScale.ScaleWithColor(color.Gray{currentGray})
	options.GeoM.Translate(currentX, currentY)
	if !b.inTransition {
		screen.DrawImage(assets.ImageMalus.SubImage(image.Rect(numBalances*gChoiceSize, 0, (numBalances+1)*gChoiceSize, gChoiceSize)).(*ebiten.Image), &options)
	}
	screen.DrawImage(assets.ImageMalus.SubImage(image.Rect(b.choices[b.choice]*gChoiceSize, 0, (b.choices[b.choice]+1)*gChoiceSize, gChoiceSize)).(*ebiten.Image), &options)
	drawLevel(screen, b.levels[b.choices[b.choice]], b.maxLevels[b.choices[b.choice]], currentX, currentY)

	// other choices
	for i := 0; i < b.numChoices-1; i++ {
		// find the choice to display
		displayNum := (b.choice + i + 1) % b.numChoices
		theChoice := b.choices[displayNum]

		//find the position to display it
		angle := float64(i+1)*(math.Pi*2)/float64(b.numChoices) + math.Pi/2 + angleShift
		x, y := math.Cos(angle)*r, -math.Sin(angle)*r
		x += float64(cX - gChoiceSize/2)
		y += float64(cY - gChoiceSize/2)

		options := ebiten.DrawImageOptions{}
		options.ColorScale.ScaleWithColor(color.Gray{gray})
		options.GeoM.Translate(x, y)
		screen.DrawImage(assets.ImageMalus.SubImage(image.Rect(theChoice*gChoiceSize, 0, (theChoice+1)*gChoiceSize, gChoiceSize)).(*ebiten.Image), &options)
		drawLevel(screen, b.levels[theChoice], b.maxLevels[theChoice], x, y)
	}

}

func (b balancing) draw(screen *ebiten.Image) {

	options := ebiten.DrawImageOptions{}
	options.GeoM.Translate(float64(gWidth-gLevelCompleteWidth)/2, float64(gTitleMargin))
	screen.DrawImage(assets.ImageLevelComplete, &options)

	b.drawChoices(screen, gWidth/2, gHeight/2-30)

	options = ebiten.DrawImageOptions{}
	options.GeoM.Translate(float64(gWidth-gTextMalusWidth)/2, float64(gHeight-gTextMalusHeight))
	id := b.choices[b.choice]
	screen.DrawImage(assets.ImageTextMalus.SubImage(image.Rect(0, id*gTextMalusHeight, gTextMalusWidth, (id+1)*gTextMalusHeight)).(*ebiten.Image), &options)
}

func newBalance(numChoices int) balancing {

	b := balancing{}

	b.choices = make([]int, numChoices)
	for i := range b.choices {
		b.choices[i] = -1
	}

	b.maxLevels[balanceGoalLines] = maxLevelGoalLines
	b.maxLevels[balanceSpeed] = maxLevelSpeed
	b.maxLevels[balanceHiddenLines] = maxLevelHiddenLines
	b.maxLevels[balanceDeathLines] = maxLevelDeathLines
	b.maxLevels[balanceInvisibleBlocks] = maxLevelInvisibleBlocks
	return b
}

func (b *balancing) getChoice() {

	possibleChoices := make([]int, 0, 2*numBalances)

BalanceLoop:
	for c := 0; c < numBalances; c++ {
		if b.levels[c] < b.maxLevels[c] {
			possibleChoices = append(possibleChoices, c)
			for _, oldChoice := range b.choices {
				if c == oldChoice {
					continue BalanceLoop
				}
			}
			possibleChoices = append(possibleChoices, c)
		}
	}

	choice := 0
	for ; len(possibleChoices) > 0 && choice < len(b.choices); choice++ {
		take := rand.Intn(len(possibleChoices))
		b.choices[choice] = possibleChoices[take]

		possibleChoices = removeElement(possibleChoices, take)

		found := true
		for found {
			found = false
			for i := range possibleChoices {
				if possibleChoices[i] == b.choices[choice] {
					possibleChoices = removeElement(possibleChoices, i)
					found = true
					break
				}
			}
		}
	}

	b.numChoices = choice

	for ; choice < len(b.choices); choice++ {
		b.choices[choice] = -1
	}

}

func (b *balancing) setChoice(choice int) {
	b.levels[choice]++
}

func (b balancing) getDeathLines() (numLines int) {
	const maxDeathLines int = 2*gPlayAreaHeightInBlocks/3 - 1

	numLines = 2*b.levels[balanceDeathLines] + 1
	if numLines > maxDeathLines {
		numLines = maxDeathLines
	}
	return
}

func (b balancing) getHiddenLines() (numLines int) {
	const hiddenFactor int = 3
	const maxHiddenLines int = 15

	numLines = hiddenFactor * b.levels[balanceHiddenLines]
	if numLines > maxHiddenLines {
		numLines = maxHiddenLines
	}

	return
}

func (b balancing) getGoalLines() int {
	var goalLines [maxLevelGoalLines + 1]int = [maxLevelGoalLines + 1]int{
		4, 8, 12,
	}

	if b.levels[balanceGoalLines] < len(goalLines) {
		return goalLines[b.levels[balanceGoalLines]]
	}
	return goalLines[len(goalLines)-1]
}

func (b balancing) getSpeedLevel(baseSpeedLevel int) int {
	var speedLevels [maxLevelSpeed]int = [maxLevelSpeed]int{
		1, 2, 4, 7, 10,
	}

	id := b.levels[balanceSpeed]
	if id >= len(speedLevels) {
		id = len(speedLevels) - 1
	}
	baseSpeedLevel += speedLevels[id]

	if baseSpeedLevel >= gSpeedLevels {
		baseSpeedLevel = gSpeedLevels - 1
	}

	return baseSpeedLevel
}

func (b balancing) getInvisibleBlocks() int {
	return b.levels[balanceInvisibleBlocks]
}
