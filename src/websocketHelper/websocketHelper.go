	package websockethelper

	import (
		"bufio"
		"context"
		"fmt"
		"io"
		"os"
		"time"
	)
	func WebsocketHelper(ctx context.Context, logPath string, callback func(string)) error {
		file, err := os.Open(logPath)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		defer file.Close()

		if _, err := file.Seek(0, io.SeekStart); err != nil {//io.SeekEnd); err != nil {
			return fmt.Errorf("seek error: %w", err)
		}

		scanner := bufio.NewScanner(file)
		scanner.Buffer(make([]byte, 64*1024), 1*1024*1024)

		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
				for scanner.Scan() {
					callback(scanner.Text())
				}

				if err := scanner.Err(); err != nil {
					return fmt.Errorf("scanner error: %w", err)
				}

				if rotated, err := isFileRotated(logPath, file); err != nil {
					return fmt.Errorf("rotation check failed: %w", err)
				} else if rotated {
					file.Close()
					if file, err = os.Open(logPath); err != nil {
						return fmt.Errorf("reopen failed: %w", err)
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