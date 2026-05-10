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
	previewImage *common.PreviewImage
	name         string
	x, y         float64
	baseScale    float64
	hovered      bool
	hoverScale   float64
}

type HomeScene struct {
	images       []*puzzleImage
	uploadButton *common.Button
	gameImage    *common.GameImage
	text         *common.TextRenderer
	startDialog  *startGameDialog
}

func NewHomeScene(gameImage *common.GameImage) *HomeScene {
	images, err := loadImages("./pictures")
	if err != nil {
		panic(err)
	}

	var puzzleImages []*puzzleImage

	screenWidth := float64(common.ScreenWidth)
	screenHeight := float64(common.ScreenHeight)

	text := common.NewTextRenderer(common.RobotoBoldFontName, common.TitleColor, 40, etxt.Center)

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

			x := leftMargin + spacingX + float64(col)*(maxImgWidthHover+spacingX)
			y := topMargin + spacingY + float64(row)*(maxImgHeightHover+spacingY) + float64(row*8)

			previewImage := common.NewPreviewImage(
				img.image, x, y, baseScale,
				common.PreviewImageOption.WithBorderColor(common.PrimaryColor),
				common.PreviewImageOption.WithCaption(text, img.name, common.BodyTextColor))
			puzzleImages = append(puzzleImages, &puzzleImage{
				previewImage: previewImage,
				name:         img.name,
				x:            x,
				y:            y,
				baseScale:    baseScale,
				hoverScale:   hoverScale,
			})
		}
	}

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
		text:        text,
		startDialog: newStartGameDialog(),
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
			h.startDialog.Open(img, name)
		}
	}

	if h.startDialog.IsOpen() {
		h.startDialog.Update()
		if h.startDialog.startClicked() {
			h.gameImage.SetPieceCount(h.startDialog.PieceCount())
			h.gameImage.SetImage(h.startDialog.ImageName(), h.startDialog.PreviewImage())
			context.SceneManager.SetScene("Game")
			return nil
		}
		if h.startDialog.cancelClicked() {
			h.startDialog.Close()
		}
		return nil
	}

	for _, img := range h.images {
		if img.previewImage.IsPointInImage(float64(mx), float64(my)) {
			img.hovered = true
		} else {
			img.hovered = false
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		for _, img := range h.images {
			if img.previewImage.IsPointInImage(float64(mx), float64(my)) {
				h.startDialog.Open(img.previewImage.Image, img.name)
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
		img.previewImage.Scale += (targetScale - img.previewImage.Scale) * 0.15

		// Card background
		cardColor := common.SurfaceColor
		if img.hovered {
			cardColor = common.SurfaceHoverColor
		}

		img.previewImage.BGColor = cardColor
		img.previewImage.Draw(screen)
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
		imgH := float64(lastImg.previewImage.Image.Bounds().Dy()) * lastImg.previewImage.Scale
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

	h.startDialog.Draw(screen)
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
