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

﻿using System;
using Microsoft.Extensions.Logging;
using NLog.Extensions.Logging;

var loggerFactory = LoggerFactory.Create(builder => {
	builder.SetMinimumLevel(Microsoft.Extensions.Logging.LogLevel.Trace);
	builder.AddNLog();
});
var logger = loggerFactory.CreateLogger("Program");

logger.LogTrace("Sample");
logger.LogDebug("Sample");
logger.LogInformation("Sample");
logger.LogWarning("Sample\n\tMultiline\n\t:)");
logger.LogError("Sample");
logger.LogCritical("Sample");
