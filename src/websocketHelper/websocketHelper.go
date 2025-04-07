package websockethelper

import (
	"bufio"
	"io"
	"main/src/logger"
	"os"
	"time"
)

func WebsocketHelper(logPath string, callback func(string), done <-chan struct{}) {
	file, err := os.Open(logPath)
	if err != nil {
		logger.Logger.Error("Failed to open log file:", err)
		return
	}
	defer file.Close()

	if _, err := file.Seek(0, io.SeekEnd); err != nil {
		logger.Logger.Error("Failed to seek log file:", err)
		return
	}

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1*1024*1024) // 1MB max line length

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			for scanner.Scan() {
				select {
				case <-done:
					return
				default:
					callback(scanner.Text())
				}
			}

			if err := scanner.Err(); err != nil {
				logger.Logger.Error("Scanner error:", err)
				return
			}

			// Проверяем не был ли файл пересоздан
			if rotated, err := isFileRotated(logPath, file); err != nil {
				logger.Logger.Error("Rotation check failed:", err)
				return
			} else if rotated {
				logger.Logger.Info("Log file rotated, reopening...")
				file.Close()
				file, err = os.Open(logPath)
				if err != nil {
					logger.Logger.Error("Failed to reopen log file:", err)
					return
				}
				scanner = bufio.NewScanner(file)
			}
		}
	}
}

func isFileRotated(path string, currentFile *os.File) (bool, error) {
	currentPos, err := currentFile.Seek(0, io.SeekCurrent)
	if err != nil {
		return false, err
	}
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.Size() < currentPos, nil
}