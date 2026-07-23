package hdcp

import (
	"fmt"
	"log"

	"github.com/thd-spatial-ai/ignis/internal/calc"
	"github.com/thd-spatial-ai/ignis/internal/models"
)

// Logger wraps the standard logger
type Logger struct {
	logger *log.Logger
}

// NewLogger creates a new Logger instance
func NewLogger(logger *log.Logger) *Logger {
	return &Logger{logger: logger}
}

// Error logs an error message
func (l *Logger) Error(format string, v ...interface{}) {
	if l.logger != nil {
		l.logger.Printf("ERROR: "+format, v...)
	}
}

// Pipeline handles the heating demand calculation pipeline, using different levels of calculations.
// Each level relies on data from the previous level. Run() stops and returns an error if any level panics.
type Pipeline struct {
	Lvl0   *models.TabulaBuildingParameters
	Logger *Logger
	err    error // set by handleError on the first panic; causes Run() to abort

	// Calculation levels
	Lvl1  *calc.CalcLevel1
	Lvl2  *calc.CalcLevel2
	Lvl3  *calc.CalcLevel3
	Lvl4  *calc.CalcLevel4
	Lvl5  *calc.CalcLevel5
	Lvl6  *calc.CalcLevel6
	Lvl7  *calc.CalcLevel7
	Lvl8  *calc.CalcLevel8
	Lvl9  *calc.CalcLevel9
	Lvl10 *calc.CalcLevel10
	Lvl11 *calc.CalcLevel11
	Lvl12 *calc.CalcLevel12
	Lvl13 *calc.CalcLevel13
	Lvl14 *calc.CalcLevel14
	Lvl15 *calc.CalcLevel15
	Lvl16 *calc.CalcLevel16
	Lvl17 *calc.CalcLevel17
}

// NewPipeline initializes the HDCP pipeline with initial calculation level and logger.
func NewPipeline(lvl0 *models.TabulaBuildingParameters, logger *Logger) *Pipeline {
	return &Pipeline{
		Lvl0:   lvl0,
		Logger: logger,
	}
}

// Run executes the full calculation pipeline from level 1 to 17.
// Returns the final heating demand (kWh/(m²·a)) and any error. If a level panics the
// pipeline stops immediately — subsequent levels are not run — and an error is returned.
func (p *Pipeline) Run() (float64, error) {
	steps := []func(){
		p.calcLvl1, p.calcLvl2, p.calcLvl3, p.calcLvl4,
		p.calcLvl5, p.calcLvl6, p.calcLvl7, p.calcLvl8,
		p.calcLvl9, p.calcLvl10, p.calcLvl11, p.calcLvl12,
		p.calcLvl13, p.calcLvl14, p.calcLvl15, p.calcLvl16,
	}
	for _, step := range steps {
		step()
		if p.err != nil {
			return 0, p.err
		}
	}
	result := p.calcLvl17()
	if p.err != nil {
		return 0, p.err
	}
	return result, nil
}

func (p *Pipeline) calcLvl1() {
	defer p.handleError("Calculation level 1")
	p.Lvl1 = calc.NewCalcLevel1(p.Lvl0)
}

func (p *Pipeline) calcLvl2() {
	defer p.handleError("Calculation level 2")
	p.Lvl2 = calc.NewCalcLevel2(p.Lvl0, p.Lvl1)
}

func (p *Pipeline) calcLvl3() {
	defer p.handleError("Calculation level 3")
	p.Lvl3 = calc.NewCalcLevel3(p.Lvl0, p.Lvl2)
}

func (p *Pipeline) calcLvl4() {
	defer p.handleError("Calculation level 4")
	p.Lvl4 = calc.NewCalcLevel4(p.Lvl0, p.Lvl1, p.Lvl2, p.Lvl3)
}

func (p *Pipeline) calcLvl5() {
	defer p.handleError("Calculation level 5")
	p.Lvl5 = calc.NewCalcLevel5(p.Lvl0, p.Lvl1, p.Lvl3, p.Lvl4)
}

func (p *Pipeline) calcLvl6() {
	defer p.handleError("Calculation level 6")
	p.Lvl6 = calc.NewCalcLevel6(p.Lvl0, p.Lvl1, p.Lvl2, p.Lvl3, p.Lvl4, p.Lvl5)
}

func (p *Pipeline) calcLvl7() {
	defer p.handleError("Calculation level 7")
	p.Lvl7 = calc.NewCalcLevel7(p.Lvl0, p.Lvl1, p.Lvl3, p.Lvl4, p.Lvl5, p.Lvl6)
}

func (p *Pipeline) calcLvl8() {
	defer p.handleError("Calculation level 8")
	p.Lvl8 = calc.NewCalcLevel8(p.Lvl0, p.Lvl1, p.Lvl2, p.Lvl4, p.Lvl5, p.Lvl6, p.Lvl7)
}

func (p *Pipeline) calcLvl9() {
	defer p.handleError("Calculation level 9")
	p.Lvl9 = calc.NewCalcLevel9(p.Lvl0, p.Lvl8)
}

func (p *Pipeline) calcLvl10() {
	defer p.handleError("Calculation level 10")
	p.Lvl10 = calc.NewCalcLevel10(p.Lvl1, p.Lvl2, p.Lvl4, p.Lvl5, p.Lvl6, p.Lvl7, p.Lvl9)
}

func (p *Pipeline) calcLvl11() {
	defer p.handleError("Calculation level 11")
	p.Lvl11 = calc.NewCalcLevel11(p.Lvl2, p.Lvl5, p.Lvl6, p.Lvl7, p.Lvl8, p.Lvl10)
}

func (p *Pipeline) calcLvl12() {
	defer p.handleError("Calculation level 12")
	p.Lvl12 = calc.NewCalcLevel12(p.Lvl0, p.Lvl1, p.Lvl11)
}

func (p *Pipeline) calcLvl13() {
	defer p.handleError("Calculation level 13")
	p.Lvl13 = calc.NewCalcLevel13(p.Lvl1, p.Lvl11, p.Lvl12)
}

func (p *Pipeline) calcLvl14() {
	defer p.handleError("Calculation level 14")
	p.Lvl14 = calc.NewCalcLevel14(p.Lvl13)
}

func (p *Pipeline) calcLvl15() {
	defer p.handleError("Calculation level 15")
	p.Lvl15 = calc.NewCalcLevel15(p.Lvl1, p.Lvl8, p.Lvl14)
}

func (p *Pipeline) calcLvl16() {
	defer p.handleError("Calculation level 16")
	p.Lvl16 = calc.NewCalcLevel16(p.Lvl13, p.Lvl15)
}

func (p *Pipeline) calcLvl17() float64 {
	defer p.handleError("Calculation level 17")
	p.Lvl17 = calc.NewCalcLevel17(p.Lvl1, p.Lvl8, p.Lvl14, p.Lvl16)
	return p.Lvl17.Run()
}

// handleError recovers from panics, logs them, and stores the error so Run() can abort.
func (p *Pipeline) handleError(level string) {
	if r := recover(); r != nil {
		p.Logger.Error("Error occurred in %s: %v", level, r)
		if p.err == nil {
			p.err = fmt.Errorf("pipeline failed at %s: %v", level, r)
		}
	}
}
