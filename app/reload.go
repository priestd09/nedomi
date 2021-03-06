package app

import (
	"errors"
	"fmt"

	"github.com/ironsmile/nedomi/config"
)

var (
	errCfgIsNil                      = errors.New("no config was provided for the reload")
	errCfgUserIsDifferent            = errors.New("can't change user by reload")
	errCfgWorkDirIsDifferent         = errors.New("can't change workdir by reload")
	errCfgListenIsDifferent          = errors.New("can't change addressed being listened to by reload")
	errCfgMaxTransferSizeIsDifferent = errors.New("can't change max_io_transfer_size by reload")
	errCfgMinTransferSizeIsDifferent = errors.New("can't change min_io_transfer_size by reload")
	errCfgReadTimeoutIsDifferent     = errors.New("can't change read_timeout by reload")
	errCfgWriteTimeoutIsDifferent    = errors.New("can't change write_timeout by reload")
	errCfgMaxHeadersSizeIsDifferent  = errors.New("can't change max_headers_size by reload")

	errTmplDifferentType      = "different types for same id '%s' between configs"
	errTmplDifferentPath      = "different paths for same id '%s' between configs"
	errTmplDifferentAlgorithm = "different algorithms for same id '%s' between configs"
	errTmplDifferentPartSize  = "different part size for same id '%s' between configs"
)

// checks if the provided config could be loaded in place of the current one.
// If that is true a nil is returned, otherwise an error explaining why it
// couldn't be done is returned
func (a *Application) checkConfigCouldBeReloaded(cfg *config.Config) error {
	if cfg == nil {
		return errCfgIsNil
	}
	if a.cfg.System.Workdir != cfg.System.Workdir {
		return errCfgWorkDirIsDifferent
	}
	if a.cfg.System.User != cfg.System.User {
		return errCfgUserIsDifferent
	}
	if a.cfg.HTTP.Listen != cfg.HTTP.Listen {
		return errCfgListenIsDifferent
	}
	if a.cfg.HTTP.MaxIOTransferSize != cfg.HTTP.MaxIOTransferSize {
		return errCfgMaxTransferSizeIsDifferent
	}
	if a.cfg.HTTP.MinIOTransferSize != cfg.HTTP.MinIOTransferSize {
		return errCfgMinTransferSizeIsDifferent
	}
	if a.cfg.HTTP.ReadTimeout != cfg.HTTP.ReadTimeout {
		return errCfgReadTimeoutIsDifferent
	}
	if a.cfg.HTTP.WriteTimeout != cfg.HTTP.WriteTimeout {
		return errCfgWriteTimeoutIsDifferent
	}
	if a.cfg.HTTP.MaxHeadersSize != cfg.HTTP.MaxHeadersSize {
		return errCfgMaxHeadersSizeIsDifferent
	}

	return cacheZonesAreCompatible(a.cfg.CacheZones, cfg.CacheZones)
}

func cacheZonesAreCompatible(zones1, zones2 map[string]*config.CacheZone) error {
	for key, zone2 := range zones2 {
		zone1 := zones1[key]
		if zone1 == nil {
			continue
		}
		if zone2.Type != zone1.Type {
			return fmt.Errorf(errTmplDifferentType, key)
		}
		if zone2.Path != zone1.Path {
			return fmt.Errorf(errTmplDifferentPath, key)
		}

		if zone2.Algorithm != zone1.Algorithm {
			return fmt.Errorf(errTmplDifferentAlgorithm, key)
		}
		if zone2.PartSize != zone1.PartSize {
			return fmt.Errorf(errTmplDifferentPartSize, key)
		}
	}
	// !TODO check that a zone does not have the same path but with different ID

	return nil
}
