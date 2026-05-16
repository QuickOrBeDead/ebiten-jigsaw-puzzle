package common

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type Scene interface {
	Update(context *SceneContext) error
	Draw(screen *ebiten.Image, context *SceneContext)
}

type SceneContext struct {
	SceneManager *SceneManager
	Cursor       ebiten.CursorShapeType
}

type SceneManager struct {
	current       Scene
	sceneContext  *SceneContext
	sceneCreators map[string]func(*SceneContext) Scene
	scenes        map[string]Scene
}

func NewSceneManager() *SceneManager {
	m := &SceneManager{
		sceneContext:  &SceneContext{},
		sceneCreators: map[string]func(*SceneContext) Scene{},
		scenes:        map[string]Scene{},
	}

	m.sceneContext.SceneManager = m
	return m
}

func (s *SceneManager) AddScene(name string, newSceneFunc func(*SceneContext) Scene) {
	s.sceneCreators[name] = newSceneFunc
}

func (s *SceneManager) SetScene(name string) {
	if scene, ok := s.scenes[name]; ok {
		s.current = scene
		return
	}

	newSceneFunc, b := s.sceneCreators[name]
	if !b {
		panic(fmt.Sprintf("%s scene not found", name))
	}

	n := newSceneFunc(s.sceneContext)
	s.scenes[name] = n
	s.current = n
}

func (s *SceneManager) Update() error {
	if s.current != nil {
		s.sceneContext.Cursor = ebiten.CursorShapeDefault
		r := s.current.Update(s.sceneContext)
		ebiten.SetCursorShape(s.sceneContext.Cursor)
		return r
	}

	return nil
}

func (s *SceneManager) Draw(screen *ebiten.Image) {
	if s.current != nil {
		s.current.Draw(screen, s.sceneContext)
	}
}
