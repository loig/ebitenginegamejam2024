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
package assets

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed squares.png
var imageSquaresBytes []byte
var ImageSquares *ebiten.Image

//go:embed back.png
var imageBackBytes []byte
var ImageBack *ebiten.Image

//go:embed digits.png
var imageDigitsBytes []byte
var ImageDigits *ebiten.Image

//go:embed malus.png
var imageMalusBytes []byte
var ImageMalus *ebiten.Image

//go:embed level.png
var imageLevelBytes []byte
var ImageLevel *ebiten.Image

//go:embed coin.png
var imageCoinBytes []byte
var ImageCoin *ebiten.Image

//go:embed bigdigits.png
var imageBigdigitsBytes []byte
var ImageBigdigits *ebiten.Image

//go:embed improvements.png
var imageImprovementsBytes []byte
var ImageImprovements *ebiten.Image

//go:embed improvementsarrow.png
var imageImprovementsArrowBytes []byte
var ImageImprovementsArrow *ebiten.Image

//go:embed max.png
var imageMaxBytes []byte
var ImageMax *ebiten.Image

//go:embed fog.png
var imageFogBytes []byte
var ImageFog *ebiten.Image

//go:embed danger.png
var imageDangerBytes []byte
var ImageDanger *ebiten.Image

//go:embed levelcomplete.png
var imageLevelCompleteBytes []byte
var ImageLevelComplete *ebiten.Image

//go:embed youlose.png
var imageYouLoseBytes []byte
var ImageYouLose *ebiten.Image

//go:embed shopfond.png
var imageShopBackBytes []byte
var ImageShopBack *ebiten.Image

//go:embed shoptitle.png
var imageShopTitleBytes []byte
var ImageShopTitle *ebiten.Image

//go:embed continue.png
var imageContinueBytes []byte
var ImageContinue *ebiten.Image

//go:embed moneyback.png
var imageMoneyBackBytes []byte
var ImageMoneyBack *ebiten.Image

//go:embed hold.png
var imageHoldBytes []byte
var ImageHold *ebiten.Image

//go:embed textesmalus.png
var imageTextMalusBytes []byte
var ImageTextMalus *ebiten.Image

//go:embed textesshop.png
var imageTextShopBytes []byte
var ImageTextShop *ebiten.Image

//go:embed heartfull.png
var imageFullHeartBytes []byte
var ImageFullHeart *ebiten.Image

//go:embed heart.png
var imageHeartBytes []byte
var ImageHeart *ebiten.Image

//go:embed controls.png
var imageControlsBytes []byte
var ImageControls *ebiten.Image

//go:embed title1.png
var imageTitle1Bytes []byte
var ImageTitle1 *ebiten.Image

//go:embed credits.png
var imageCreditsBytes []byte
var ImageCredits *ebiten.Image

//go:embed title2.png
var imageTitle2Bytes []byte
var ImageTitle2 *ebiten.Image

//go:embed win.png
var imageWinBytes []byte
var ImageWin *ebiten.Image

//go:embed rocket.png
var imageRocketBytes []byte
var ImageRocket *ebiten.Image

func Load(mult int) {
	var err error

	imageDecoded, _, err := image.Decode(bytes.NewReader(imageSquaresBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageSquares = ebiten.NewImageFromImage(imageDecoded)
	// resize
	ImageSquares = resize(ImageSquares, mult)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageBackBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageBack = ebiten.NewImageFromImage(imageDecoded)
	// resize
	ImageBack = resize(ImageBack, mult)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageDigitsBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageDigits = ebiten.NewImageFromImage(imageDecoded)
	// resize
	ImageDigits = resize(ImageDigits, mult)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageMalusBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageMalus = ebiten.NewImageFromImage(imageDecoded)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageLevelBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageLevel = ebiten.NewImageFromImage(imageDecoded)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageCoinBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageCoin = ebiten.NewImageFromImage(imageDecoded)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageBigdigitsBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageBigdigits = ebiten.NewImageFromImage(imageDecoded)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageImprovementsBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageImprovements = ebiten.NewImageFromImage(imageDecoded)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageImprovementsArrowBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageImprovementsArrow = ebiten.NewImageFromImage(imageDecoded)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageMaxBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageMax = ebiten.NewImageFromImage(imageDecoded)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageFogBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageFog = ebiten.NewImageFromImage(imageDecoded)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageDangerBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageDanger = ebiten.NewImageFromImage(imageDecoded)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageLevelCompleteBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageLevelComplete = ebiten.NewImageFromImage(imageDecoded)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageYouLoseBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageYouLose = ebiten.NewImageFromImage(imageDecoded)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageShopBackBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageShopBack = ebiten.NewImageFromImage(imageDecoded)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageShopTitleBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageShopTitle = ebiten.NewImageFromImage(imageDecoded)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageContinueBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageContinue = ebiten.NewImageFromImage(imageDecoded)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageMoneyBackBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageMoneyBack = ebiten.NewImageFromImage(imageDecoded)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageHoldBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageHold = ebiten.NewImageFromImage(imageDecoded)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageTextMalusBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageTextMalus = ebiten.NewImageFromImage(imageDecoded)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageTextShopBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageTextShop = ebiten.NewImageFromImage(imageDecoded)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageFullHeartBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageFullHeart = ebiten.NewImageFromImage(imageDecoded)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageHeartBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageHeart = ebiten.NewImageFromImage(imageDecoded)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageControlsBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageControls = ebiten.NewImageFromImage(imageDecoded)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageTitle1Bytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageTitle1 = ebiten.NewImageFromImage(imageDecoded)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageCreditsBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageCredits = ebiten.NewImageFromImage(imageDecoded)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageTitle2Bytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageTitle2 = ebiten.NewImageFromImage(imageDecoded)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageWinBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageWin = ebiten.NewImageFromImage(imageDecoded)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageRocketBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageRocket = ebiten.NewImageFromImage(imageDecoded)
}

func resize(img *ebiten.Image, mult int) (res *ebiten.Image) {

	res = ebiten.NewImage(img.Bounds().Dx()*mult, img.Bounds().Dy()*mult)

	options := ebiten.DrawImageOptions{}
	options.GeoM.Scale(float64(mult), float64(mult))

	res.DrawImage(img, &options)

	return

}
