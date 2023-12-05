#ifndef SETTINGS_H
#define SETTINGS_H

#include <QString>
#include <QStringList>
#include <QSettings>
#include <QDebug>
#include <QList>
#include <QMap>
#include <QMessageBox>
#include <QUuid>
#include <QFileInfo>
#include <QDir>
#include <QMutex>
#include <QDateTime>
#include <QFile>
#include <QTextStream>

void log_handler(QtMsgType type, const QMessageLogContext &info, const QString &msg);

class Settings : public QSettings
{
public:
    static Settings* GetInstance();
    void outputLog(const QString &type, const char* file, const char* func, int line, const QString &msg);
    ~Settings();

private:
    Settings();

public:
    static QString APP_NAME;

private:
    static Settings* m_iInstance;
    QMutex          m_mutex;
    QString m_LogPath;
    QDateTime m_LastLogTime;
    QFile* m_iLogFile;
    QTextStream* m_iLogStream;
};

#endif // SETTINGS_H


