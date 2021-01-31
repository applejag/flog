using System;
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
logger.LogWarning("Sample");
logger.LogError("Sample");
logger.LogCritical("Sample");

