package installer

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Installer struct {
	srcDir           string
	templatePath     string
	skipPatterns     []string
	concurrencyLimit int
	templateHeader   string // (Шаблон)
}

func New(
	srcDir,
	templatePath string,
	skipPatterns []string,
	concurrencyLimit int,
) *Installer {
	return &Installer{
		srcDir:           srcDir,
		templatePath:     templatePath,
		skipPatterns:     skipPatterns,
		concurrencyLimit: concurrencyLimit,
	}
}

func (in *Installer) Run() error {
	// 1 - Шаблон
	tplBytes, err := os.ReadFile(in.templatePath)
	if err != nil {
		return fmt.Errorf("не удалось прочитать шаблон %q: %v", in.templatePath, err)
	}
	in.templateHeader = string(tplBytes)

	// 2 Многопоточность
	var wg sync.WaitGroup
	sem := make(chan struct{}, in.concurrencyLimit) // Семафор

	// 3 Запуск
	wg.Add(1)
	go in.walkDirConcurrently(in.srcDir, &wg, sem)
	wg.Wait()

	return nil
}

func (in *Installer) walkDirConcurrently(dirPath string, wg *sync.WaitGroup, sem chan struct{}) {
	defer wg.Done()

	sem <- struct{}{}
	defer func() { <-sem }()

	if err := in.walkDir(dirPath, wg, sem); err != nil {
		log.WithError(err).Warnf("ошибка обхода директории %v", in.srcDir)
	}
}

func (in *Installer) walkDir(dirPath string, wg *sync.WaitGroup, sem chan struct{}) error {
	dirName := filepath.Base(dirPath)
	if in.shouldSkipDir(dirName) {
		log.Infof("Пропускаем директорию %v из-за паттернов", dirName)
		return nil
	}

	entries, err := os.ReadDir(dirPath) //entries
	if err != nil {
		return fmt.Errorf("ошибка чтения директории %q: %v", dirPath, err)
	}

	for _, entry := range entries {
		name := entry.Name()
		fullPath := filepath.Join(dirPath, name)

		if entry.IsDir() {
			wg.Add(1)
			go in.walkDirConcurrently(dirPath, wg, sem)
		} else {
			if filepath.Ext(name) == ".md" {
				if err := processFile(fullPath); err != nil {
					log.WithError(err).Warnf("Ошибка обработка файла: %s", fullPath)
				}
			}
		}
	}

	return nil
}

func processFile(path string) error {
	return nil
}

// Проверяет, нужно ли пропускать директорию по паттерну
func (in *Installer) shouldSkipDir(dirName string) bool {
	for _, pattern := range in.skipPatterns {
		if strings.HasPrefix(dirName, pattern) {
			return true
		}
	}

	return false
}
