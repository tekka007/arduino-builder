/*
 * This file is part of Arduino Builder.
 *
 * Arduino Builder is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
 *
 * As a special exception, you may use this file as part of a free software
 * library without restriction.  Specifically, if other files instantiate
 * templates or use macros or inline functions from this file, or you compile
 * this file and link it with other files to produce an executable, this
 * file does not by itself cause the resulting executable to be covered by
 * the GNU General Public License.  This exception does not however
 * invalidate any other reasons why the executable file might be covered by
 * the GNU General Public License.
 *
 * Copyright 2015 Arduino LLC (http://www.arduino.cc/)
 */

package phases

import (
	"arduino.cc/builder/builder_utils"
	"arduino.cc/builder/constants"
	"arduino.cc/builder/i18n"
	"arduino.cc/builder/types"
	"arduino.cc/builder/utils"
	"os"
	"path/filepath"
)

type SketchBuilder struct{}

func (s *SketchBuilder) Run(ctx *types.Context) error {
	sketchBuildPath := ctx.SketchBuildPath
	buildProperties := ctx.BuildProperties
	includes := ctx.IncludeFolders
	includes = utils.Map(includes, utils.WrapWithHyphenI)
	verbose := ctx.Verbose
	warningsLevel := ctx.WarningsLevel
	logger := ctx.GetLogger()

	err := utils.EnsureFolderExists(sketchBuildPath)
	if err != nil {
		return i18n.WrapError(err)
	}

	var objectFiles []string
	objectFiles, err = builder_utils.CompileFiles(objectFiles, sketchBuildPath, false, sketchBuildPath, buildProperties, includes, verbose, warningsLevel, logger)
	if err != nil {
		return i18n.WrapError(err)
	}

	// The "src/" subdirectory of a sketch is compiled recursively
	sketchSrcPath := filepath.Join(sketchBuildPath, constants.SKETCH_FOLDER_SRC)
	if info, err := os.Stat(sketchSrcPath); err == nil && info.IsDir() {
		objectFiles, err = builder_utils.CompileFiles(objectFiles, sketchSrcPath, true, sketchSrcPath, buildProperties, includes, verbose, warningsLevel, logger)
		if err != nil {
			return i18n.WrapError(err)
		}
	}

	ctx.SketchObjectFiles = objectFiles

	return nil
}
