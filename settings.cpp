#include "settings.h"
#include <QFile>
#include <QJsonParseError>
#include <QDebug>
#include <QJsonDocument>
#include <QJsonObject>
#include <QJsonArray>
#include <QDirIterator>
#include <QRandomGenerator>
#include "go_interface.h"



void log(QtMsgType type, const QMessageLogContext &info, const QString &msg)
{
    QString lType;
    switch (type) {
    case QtDebugMsg:
        lType = "D";
        break;
    case QtInfoMsg:
        lType = "I";
        break;
    case QtWarningMsg:
        lType = "W";
        break;
    case QtCriticalMsg:
        lType = "E";
        break;
    case QtFatalMsg:
        lType = "F";
        break;
    default:
        break;
    }
    Settings::m_iInstance->outputLog(lType, info.file, info.function, info.line, msg);
}

static bool readJsonFile(QIODevice &device, QSettings::SettingsMap &map)
{
    QJsonParseError jsonError;
    QJsonDocument jsonDoc(QJsonDocument::fromJson(device.readAll(), &jsonError));
    if(jsonError.error != QJsonParseError::NoError)
    {
        qDebug()<< "json error:" << jsonError.errorString();
        return false;
    }
    map = jsonDoc.toVariant().toMap();
    //    for(QMap<QString, QVariant>::const_iterator iter = map1.begin(); iter != map1.end(); ++iter){
    //        map[iter.key()] = iter.value();
    //    }
    return true;
}

static bool writeJsonFile(QIODevice &device, const QSettings::SettingsMap &map)
{
    bool ret = false;

    QJsonObject rootObj;
    /*QJsonDocument jsonDocument; = QJsonDocument::fromVariant(QVariant::fromValue(map));
    if ( device.write(jsonDocument.toJson()) != -1 )
        ret = true;*/
    for(QMap<QString, QVariant>::const_iterator iter = map.begin(); iter != map.end(); ++iter){
        rootObj[iter.key()] = QJsonValue::fromVariant(iter.value());
    }
    QJsonDocument jsonDoc;
    jsonDoc.setObject(rootObj);
    if ( device.write(jsonDoc.toJson()) != -1 )
        ret = true;
    return ret;
}

Settings* Settings::m_iInstance = nullptr;
QString Settings::APP_NAME = "models_gen";

Settings* Settings::GetInstance(){
    if(m_iInstance == nullptr){
        m_iInstance = new Settings();
        qInstallMessageHandler(log);
    }
    return m_iInstance;
}



Settings::Settings():QSettings(Settings::APP_NAME+".json" ,QSettings::registerFormat("json", readJsonFile, writeJsonFile))
{
    sync();
    m_LogPath = "logs";

    pbapi::PK_SET_LOG_DIR_REQ req;
    req.set_dir("logs");
    pbapi::PK_SET_LOG_DIR_RSP rsp;
    Call_Go_Func(&req, &rsp, SetLogDir);
}

Settings::~Settings()
{
    pbapi::PK_SET_LOG_DIR_REQ req;
    pbapi::PK_SET_LOG_DIR_RSP rsp;
    Call_Go_Func(&req, &rsp, SetLogDir);
    printf("-----------Settings destruction\n");
}

void Settings::outputLog(const QString &type, const char* file, const char* func, int line, const QString &msg)
{
    QMutexLocker locker(&m_mutex);
    QString time = QDateTime::currentDateTime().toString("yyyy-MM-dd hh:mm:ss.zzz");
    QString str = QString("[%1][%2 %3:%5][%4] ===> %6\n")
                      .arg(type).arg(file).arg(func).arg(time).arg(line).arg(msg);
    QString name = QString("%1/%2.txt").arg(m_LogPath).arg(QDateTime::currentDateTime().toString("yyyy-MM-dd"));
    QDir dir(m_LogPath);
    if(!dir.exists())
        dir.mkpath(m_LogPath);
    QFile           m_file(name);
    if(!m_file.open(QIODevice::WriteOnly | QIODevice::Append))
        return;
    m_file.write(str.toUtf8());
    m_file.close();
}

