export interface LogEntry {
    level: string,
    message: string,
    data: any,
    timestamp: string,
}

export enum LogLevel {
    FATAL = "fatal",
    PANIC = "panic",
    ERROR = "error",
    WARN = "warn",
    INFO = "info",
    DEBUG = "debug",
    TRACE = "trace",
}

export function getLogLevel(level: number): string {
    switch (level) {
        case 1:
            return LogLevel.FATAL;
        case 2:
            return LogLevel.PANIC;
        case 3:
            return LogLevel.ERROR;
        case 4:
            return LogLevel.WARN;
        case 5:
            return LogLevel.INFO;
        case 6:
            return LogLevel.DEBUG;
        case 7:
            return LogLevel.TRACE;
    }
}
