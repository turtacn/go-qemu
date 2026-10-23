package dryrun

import (
	"context"
)

func Run(ctx context.Context, cfg Config) (*Result, error) {
	result := &Result{
		OK:     true,
		Layers: make(map[LayerMask]*LayerResult),
	}

	if cfg.Layers&Layer1Static != 0 {
		lr, err := runLayer1(ctx, cfg)
		if err != nil {
			return nil, err
		}
		result.Layers[Layer1Static] = lr
		if !lr.OK {
			result.OK = false
			return result, nil
		}
	}

	if cfg.Layers&Layer2Smoke != 0 {
		lr, err := runLayer2(ctx, cfg)
		if err != nil {
			return nil, err
		}
		result.Layers[Layer2Smoke] = lr
		if !lr.OK {
			result.OK = false
			return result, nil
		}
	}

	if cfg.Layers&Layer3Runtime != 0 {
		lr, err := runLayer3(ctx, cfg)
		if err != nil {
			return nil, err
		}
		result.Layers[Layer3Runtime] = lr
		if !lr.OK {
			result.OK = false
		}
	}

	return result, nil
}
