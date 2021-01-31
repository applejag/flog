using System;
using Microsoft.Extensions.Logging;

var loggerFactory = new LoggerFactory();
loggerFactory.AddConsole();
var logger = loggerFactory.CreateLogger<Program>();

logger.LogTrace("Sample");
logger.LogDebug("Sample");
logger.LogInformation("Sample");
logger.LogWarning("Sample");
logger.LogError("Sample");
logger.LogCritical("Sample");

