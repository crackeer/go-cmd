package util

import (
	"fmt"
	"os"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

// ProgressBarKeeper
type ProgressBarKeeper struct {
	bar *progressbar.ProgressBar
}

// DefaultProgressBarKeeper
//
//	@param total
//	@param message
//	@return *ProgressBarKeeper
//	@return error
func DefaultProgressBarKeeper(total int) (*ProgressBarKeeper, error) {
	bar := progressbar.NewOptions(total, progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(20),
		progressbar.OptionShowCount(),
		progressbar.OptionThrottle(615*time.Millisecond),
		progressbar.OptionSpinnerCustom([]string{"/", "\\"}),
		progressbar.OptionShowDescriptionAtLineEnd(),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionSetRenderBlankState(false),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		//progressbar.OptionSetDescription(formatProgressBarDescription(total, 1, message)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
	return &ProgressBarKeeper{
		bar: bar,
	}, nil
}

// IncrProgressBy ...
//
//	@receiver keeper
//	@param step
//	@param message
func (keeper *ProgressBarKeeper) IncrProgressBy(step int) {

	keeper.bar.Add(step)
}

func (keeper *ProgressBarKeeper) Describe(message string) {
	keeper.bar.Describe(message)
}

// UpdateProgress
//
//	@receiver keeper
//	@param step
//	@param message
func (keeper *ProgressBarKeeper) UpdateProgress(progress int, message string) {
	keeper.bar.Describe(message)
	keeper.bar.Set(progress)
}
