// Copyright (C) 2021  Kalle Jillheden
// SPDX-FileCopyrightText: 2021 Kalle Fagerberg
//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	_ "embed"
	"fmt"
	"os"
)

func showLicenseWarrantyThenExit() {
	fmt.Println(licenseWarranty)
	os.Exit(0)
}

func showLicenseConditionsThenExit() {
	fmt.Println(licenseConditions)
	os.Exit(0)
}

const LicenceNotice = `flog  Copyright (C) 2021  Kalle Jillheden
    License GPLv3+: GNU GPL version 3 or later <https://gnu.org/licenses/gpl.html>
    This program comes with ABSOLUTELY NO WARRANTY; for details run 'flog --license-w'.
    This is free software, and you are welcome to redistribute it
    under certain conditions; run 'flog --license-c' for details.`

//go:embed assets/GPL-3.0-or-later-warranty.txt
var licenseWarranty string

//go:embed assets/GPL-3.0-or-later-conditions.txt
var licenseConditions string
