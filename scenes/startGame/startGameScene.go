package startGame

import (
	"image/color"
	"image/jpeg"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
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

type StartGameScene struct {
	images       []*puzzleImage
	uploadButton *common.Button
	gameImage    *common.GameImage
	text         *common.TextRenderer
	startDialog  *startGameDialog
	backBtn      *common.Button
}

func NewStartGameScene(gameImage *common.GameImage, context *common.SceneContext) *StartGameScene {
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

		topMargin := 120.0
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

	return &StartGameScene{
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
		backBtn: common.NewButton(
			20, 12,
			80, 40,
			"Back",
			common.ButtonOption.WithFontSize(18),
			common.ButtonOption.WithFontColor(common.BodyTextColor),
			common.ButtonOption.WithColor(common.HeaderButtonColor),
			common.ButtonOption.WithHoverColor(common.HeaderButtonHoverColor),
			common.ButtonOption.WithOnClick(func() {
				context.SceneManager.SetScene("home")
			}),
		),
		text:        text,
		startDialog: newStartGameDialog(),
	}
}

func (s *StartGameScene) Update(context *common.SceneContext) error {
	mx, my := ebiten.CursorPosition()

	s.backBtn.Update(context)
	s.uploadButton.Update(context)

	if s.uploadButton.Clicked {
		name, img, err := loadImageFromDesktop()
		if err != nil && err != dialog.Cancelled {
			return err
		}

		if img != nil {
			s.startDialog.Open(img, name)
		}
	}

	if s.startDialog.IsOpen() {
		s.startDialog.Update(context)
		if s.startDialog.startClicked() {
			s.gameImage.SetPieceCount(s.startDialog.PieceCount())
			s.gameImage.SetImage(s.startDialog.ImageName(), s.startDialog.PreviewImage())
			context.SceneManager.SetScene("game")
			return nil
		}
		if s.startDialog.cancelClicked() {
			s.startDialog.Close()
		}
		return nil
	}

	for _, img := range s.images {
		if img.previewImage.IsPointInImage(float64(mx), float64(my)) {
			img.hovered = true
			context.Cursor = ebiten.CursorShapeCrosshair
		} else {
			img.hovered = false
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		for _, img := range s.images {
			if img.previewImage.IsPointInImage(float64(mx), float64(my)) {
				s.startDialog.Open(img.previewImage.Image, img.name)
			}
		}
	}

	return nil
}

func (s *StartGameScene) Draw(screen *ebiten.Image, context *common.SceneContext) {
	screen.Fill(common.BackgroundColor)

	const headerH float32 = 120
	common.DrawPanel(screen, 0, 0, float32(common.ScreenWidth), headerH, common.HeaderColor, false, nil)

	vector.FillRect(screen, 0, 0, float32(common.ScreenWidth), 1, color.RGBA{255, 255, 255, 14}, false)

	common.DrawPanel(screen, 0, headerH-3, float32(common.ScreenWidth), 3, common.PrimaryColor, false, nil)

	for i := 0; i < 4; i++ {
		alpha := uint8(35 - 7*i)
		vector.FillRect(screen, 0, headerH+float32(i), float32(common.ScreenWidth), 1, color.RGBA{0, 0, 0, alpha}, false)
	}

	s.text.SetColor(common.TitleColor)
	s.text.SetSize(42)
	s.text.DrawEmbossedAutoWithShadow(screen, "Select a Puzzle Image", common.ScreenWidth/2, 50, color.RGBA{0, 0, 0, 100}, 2, 2)

	s.backBtn.Draw(screen)

	for _, img := range s.images {
		targetScale := img.baseScale
		if img.hovered {
			targetScale = img.hoverScale
		}
		img.previewImage.Scale += (targetScale - img.previewImage.Scale) * 0.15

		cardColor := common.SurfaceColor
		if img.hovered {
			cardColor = common.SurfaceHoverColor
		}

		img.previewImage.BGColor = cardColor
		img.previewImage.Draw(screen)
	}

	if len(s.images) > 0 {
		lastImg := s.images[len(s.images)-1]
		imgH := float64(lastImg.previewImage.ScaledH)
		orY := int(lastImg.y + imgH + 70)
		s.text.SetColor(common.MutedTextColor)
		s.text.SetSize(18)
		s.text.DrawEmbossedHozCenterWithShadow(screen, "OR", orY, color.RGBA{0, 0, 0, 80}, 1, 1)
	} else {
		s.text.SetColor(common.MutedTextColor)
		s.text.SetSize(18)
		s.text.DrawEmbossedHozCenterWithShadow(screen, "OR", 300, color.RGBA{0, 0, 0, 80}, 1, 1)
	}

	s.uploadButton.Draw(screen)

	s.startDialog.Draw(screen)
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
