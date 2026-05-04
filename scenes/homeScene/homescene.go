package homeScene

import (
	"image/color"
	"image/jpeg"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sqweek/dialog"
	"github.com/tinne26/etxt"

	"github.com/QuickOrBeDead/ebiten-jigsaw-puzzle/common"
)

type imageWithName struct {
	image *ebiten.Image
	name  string
}

type puzzleImage struct {
	image      *ebiten.Image
	name       string
	x, y       float64
	scale      float64
	baseScale  float64
	hovered    bool
	hoverScale float64
}

type HomeScene struct {
	images       []*puzzleImage
	uploadButton *common.Button
	gameImage    *common.GameImage
	text         *common.TextRenderer
}

func NewHomeScene(gameImage *common.GameImage) *HomeScene {
	images, err := loadImages("./pictures")
	if err != nil {
		panic(err)
	}

	var puzzleImages []*puzzleImage

	screenWidth := float64(common.ScreenWidth)
	screenHeight := float64(common.ScreenHeight)

	baseScale := 0.32
	hoverScale := baseScale * 1.05

	numImages := len(images)

	if numImages > 0 {
		numCols := 5
		if numCols > numImages {
			numCols = numImages
		}

		numRows := (numImages + numCols - 1) / numCols

		topMargin := 160.0
		bottomMargin := 200.0
		leftMargin := 60.0
		rightMargin := 60.0

		availableWidth := screenWidth - leftMargin - rightMargin
		availableHeight := screenHeight - topMargin - bottomMargin

		maxImgWidthHover := 0.0
		maxImgHeightHover := 0.0
		for _, img := range images {
			w := float64(img.image.Bounds().Dx()) * hoverScale
			h := float64(img.image.Bounds().Dy()) * hoverScale
			if w > maxImgWidthHover {
				maxImgWidthHover = w
			}
			if h > maxImgHeightHover {
				maxImgHeightHover = h
			}
		}

		spacingX := (availableWidth - float64(numCols)*maxImgWidthHover) / float64(numCols+1)
		spacingY := (availableHeight - float64(numRows)*maxImgHeightHover) / float64(numRows+1)

		if spacingX < 10 {
			spacingX = 10
		}
		if spacingY < 10 {
			spacingY = 10
		}

		for i, img := range images {
			col := i % numCols
			row := i / numCols

			x := leftMargin + spacingX + float64(col)*(maxImgWidthHover+spacingX) + maxImgWidthHover/2
			y := topMargin + spacingY + float64(row)*(maxImgHeightHover+spacingY) + maxImgHeightHover/2 + float64(row*8)

			puzzleImages = append(puzzleImages, &puzzleImage{
				image:      img.image,
				name:       img.name,
				x:          x,
				y:          y,
				scale:      baseScale,
				baseScale:  baseScale,
				hoverScale: hoverScale,
			})
		}
	}

	text := common.NewTextRenderer(common.RobotoBoldFontName, common.TitleColor, 40, etxt.Center)

	buttonY := screenHeight - 120.0

	return &HomeScene{
		images:    puzzleImages,
		gameImage: gameImage,
		uploadButton: common.NewButton(
			float32((screenWidth-(220+48))/2), float32(buttonY),
			220, 48,
			"Upload Image",
			common.ButtonOption.WithFontSize(22),
			common.ButtonOption.WithFontColor(common.BodyTextColor),
			common.ButtonOption.WithColor(common.PrimaryColor),
			common.ButtonOption.WithHoverColor(common.PrimaryHoverColor),
			common.ButtonOption.WithShadowColor(common.ShadowColor),
		),
		text: text,
	}
}

func (h *HomeScene) Update(context *common.SceneContext) error {
	mx, my := ebiten.CursorPosition()

	h.uploadButton.Update()

	if h.uploadButton.Clicked {
		name, img, err := loadImageFromDesktop()
		if err != nil && err != dialog.Cancelled {
			return err
		}

		if img != nil {
			h.gameImage.SetImage(name, img)
			context.SceneManager.SetScene("Game")
			return nil
		}
	}

	for _, img := range h.images {
		if isPointInImage(float64(mx), float64(my), img) {
			img.hovered = true
		} else {
			img.hovered = false
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		for _, img := range h.images {
			if isPointInImage(float64(mx), float64(my), img) {
				h.gameImage.SetImage(img.name, img.image)
				context.SceneManager.SetScene("Game")
			}
		}
	}

	return nil
}

func (h *HomeScene) Draw(screen *ebiten.Image, context *common.SceneContext) {
	// Background
	screen.Fill(common.BackgroundColor)

	// Header area
	common.DrawPanel(screen, 0, 0, float32(common.ScreenWidth), 120, common.HeaderColor, false, nil)

	// Title with shadow
	h.text.SetColor(common.SuccessColor)
	h.text.SetSize(42)
	h.text.DrawWithShadow(screen, "Welcome to Jigsaw Puzzle!", common.ScreenWidth/2, 50, color.RGBA{0, 0, 0, 100}, 2, 2)

	// Subtitle
	h.text.SetColor(common.MutedTextColor)
	h.text.SetSize(20)
	h.text.DrawHorizontalCenter(screen, "Select a puzzle image", 100)

	// Draw image cards
	isHovered := false
	for _, img := range h.images {
		// Smooth hover animation
		targetScale := img.baseScale
		if img.hovered {
			targetScale = img.hoverScale
			isHovered = true
		}
		img.scale += (targetScale - img.scale) * 0.15

		imgW := float64(img.image.Bounds().Dx())
		imgH := float64(img.image.Bounds().Dy())
		scaledW := imgW * img.scale
		scaledH := imgH * img.scale

		// Draw card background (centered)
		cardX := float32(img.x - scaledW/2)
		cardY := float32(img.y - scaledH/2)
		cardWidth := float32(scaledW)
		cardHeight := float32(scaledH)

		// Card shadow
		common.DrawSoftShadow(screen, cardX, cardY, cardWidth, cardHeight,
			color.RGBA{0, 0, 0, 80}, 4, 4)

		// Card background
		cardColor := common.SurfaceColor
		if img.hovered {
			cardColor = common.SurfaceHoverColor
		}
		common.DrawPanel(screen, cardX-5, cardY-5, cardWidth+10, cardHeight+10, cardColor, true, common.PrimaryColor)

		// Draw image (centered transform to prevent hover overlap)
		geoM := ebiten.GeoM{}
		geoM.Scale(img.scale, img.scale)
		geoM.Translate(img.x-imgW*img.scale/2, img.y-imgH*img.scale/2)
		opt := &ebiten.DrawImageOptions{GeoM: geoM, Filter: ebiten.FilterLinear}
		screen.DrawImage(img.image, opt)

		// Image name caption
		h.text.SetColor(common.BodyTextColor)
		h.text.SetSize(14)
		h.text.SetAlign(etxt.Center)
		h.text.Draw(screen, img.name, int(img.x), int(img.y+scaledH/2)+20)
	}

	// Cursor update
	if isHovered {
		ebiten.SetCursorShape(ebiten.CursorShapeCrosshair)
	} else {
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}

	// Upload section - position "OR" text between images and button
	if len(h.images) > 0 {
		lastImg := h.images[len(h.images)-1]
		imgH := float64(lastImg.image.Bounds().Dy()) * lastImg.scale
		orY := int(lastImg.y + imgH/2 + 60)
		h.text.SetColor(common.MutedTextColor)
		h.text.SetSize(18)
		h.text.DrawHorizontalCenter(screen, "OR", orY)
	} else {
		h.text.SetColor(common.MutedTextColor)
		h.text.SetSize(18)
		h.text.DrawHorizontalCenter(screen, "OR", 300)
	}

	h.uploadButton.Draw(screen)
}

func isPointInImage(x, y float64, puzzleImage *puzzleImage) bool {
	imgW := float64(puzzleImage.image.Bounds().Dx())
	imgH := float64(puzzleImage.image.Bounds().Dy())

	geoM := ebiten.GeoM{}
	geoM.Scale(puzzleImage.scale, puzzleImage.scale)
	geoM.Translate(puzzleImage.x-imgW*puzzleImage.scale/2, puzzleImage.y-imgH*puzzleImage.scale/2)

	if !geoM.IsInvertible() {
		return false
	}

	geoM.Invert()

	imgX, imgY := geoM.Apply(x, y)

	w, h := puzzleImage.image.Bounds().Dx(), puzzleImage.image.Bounds().Dy()
	return imgX >= 0 && imgX < float64(w) && imgY >= 0 && imgY < float64(h)
}

func loadImageFromDesktop() (string, *ebiten.Image, error) {
	path, err := dialog.File().Filter("Image files", "jpg", "jpeg").Load()
	if err != nil {
		return "", nil, err
	}

	img, err := loadJpegImageFromPath(path)

	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)), img, err
}

func loadImages(path string) ([]*imageWithName, error) {
	var images []*imageWithName

	err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			if ext := filepath.Ext(p); ext == ".jpg" || ext == ".jpeg" {
				img, err := loadJpegImageFromPath(p)
				if err != nil {
					return err
				}

				// get the file name without the path and extension
				name := strings.TrimSuffix(filepath.Base(p), ext)

				images = append(images, &imageWithName{
					image: img,
					name:  name,
				})
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return images, nil
}

func loadJpegImageFromPath(path string) (*ebiten.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, err
	}

	return ebiten.NewImageFromImage(img), nil
}
