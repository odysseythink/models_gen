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



class Settings : public QSettings
{
public:
    static Settings* m_iInstance;
    static Settings* GetInstance();
    void outputLog(const QString &type, const char* file, const char* func, int line, const QString &msg);
    ~Settings();

private:
    Settings();

public:
    static QString APP_NAME;

private:
    QMutex          m_mutex;
    QString m_LogPath;
};

#endif // SETTINGS_H


