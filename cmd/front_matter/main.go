package main

import (
	"flag"
	"github.com/ANkulagin/golang_yaml_front_matter_installer_sb/internal/config"
	"github.com/ANkulagin/golang_yaml_front_matter_installer_sb/internal/service/installer"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func main() {
	configPath := flag.String("config", "configs/config.yaml", "Путь к конфигурационному файлу")
	flag.Parse()
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Не удалось загрузить конфигурацию: %v", err)
	}

	// Настройка уровня логирования
	level, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatalf("Не удалось установить уровень логирования: %v", err)
	}
	log.SetLevel(level)

	// Настройка формата логирования
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		// Можно добавить другие настройки
		ForceColors: true,
	})

	// Настройка вывода логов (можно перенаправить в файл, если нужно)
	log.SetOutput(os.Stdout)

	// Преобразование относительных путей в абсолютные
	absSrcDir, err := filepath.Abs(cfg.SrcDir)
	if err != nil {
		log.Fatalf("Не удалось определить абсолютный путь для исходной директории: %v", err)
	}

	absTemplateDir, err := filepath.Abs(cfg.TemplateDir)
	if err != nil {
		log.Fatalf("Не удалось определить абсолютный путь для директории шаблонов: %v", err)
	}

	log.Infof("Уровень логирования: %s", cfg.LogLevel)
	log.Infof("Абсолютный путь к исходной директории: %s", absSrcDir)
	log.Infof("Абсолютный путь к директории шаблонов: %s", absTemplateDir)

	in := installer.New(absSrcDir, absTemplateDir, cfg.SkipPatterns, cfg.ConcurrencyLimit)
	if err := in.Run(); err != nil {
		log.Fatalf("Ошибка при выполнении установки front matter: %v", err)
	}

	log.Info("Успешное завершение работы")
}
