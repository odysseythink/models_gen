#ifndef __GO_INTERFACE_H__
#define __GO_INTERFACE_H__


#include "google/protobuf/descriptor.h"
#include "google/protobuf/message.h"
#include "api.pb.h"
#include "go_models_gen/c_2_go_interface.h"
#include "go_models_gen/libmodels_gen.h"
#include <QObject>
#include <QMap>
#include <QList>
#include <QMutex>

typedef void* (*GO_Interface_Func_Type)(void* req, int reqlen, int* l);
int Call_Go_Func(google::protobuf::Message* req, google::protobuf::Message* rsp, GO_Interface_Func_Type callback);

class Callback : public QObject
{
    Q_OBJECT
public:
    explicit Callback(int cmd, const QObject *receiver, const char *method, QObject *parent = Q_NULLPTR):QObject(parent)
    {
        m_nCmd = cmd;
        m_iReceiver = receiver;
        m_iMethod = method;
    }
    ~Callback(){

    }
    const QObject *Receiver(){
        return m_iReceiver;
    }
    const char *Method(){
        return m_iMethod;
    }

private:

signals:
    void sigActive(QByteArray msg);

private:
    int m_nCmd;
    const QObject *m_iReceiver;
    const char *m_iMethod;
};

class CGoInterface : public QObject
{
    Q_OBJECT
public:
    static CGoInterface* GetInstance(){
        if(m_iInstance == nullptr){
            m_iInstance = new CGoInterface();
        }
        return m_iInstance;
    }
    void RegisterHandler(int cmd, const QObject *receiver, const char *method){
        m_CallbacksMux.lock();
        if (!m_Callbacks.contains(cmd)){
            QList<Callback *> tmp;
            m_Callbacks[cmd] = tmp;
        }
        Callback * cb = new Callback(cmd, receiver, method);
        m_Callbacks[cmd].append(cb);
        connect(cb, SIGNAL(sigActive(QByteArray)), receiver, method);
        m_CallbacksMux.unlock();
    }

    void UnRegisterHandler(int cmd, const QObject *receiver){
        m_CallbacksMux.lock();
        if (!m_Callbacks.contains(cmd)){
            QList<Callback *> tmp;
            m_Callbacks[cmd] = tmp;
        }
        QList<Callback *>::iterator iter = m_Callbacks[cmd].begin();
        while(iter != m_Callbacks[cmd].end()){
            Callback * cb = *iter;
            if (cb->Receiver() == receiver){
                m_Callbacks[cmd].erase(iter);
                disconnect(cb, SIGNAL(sigActive(QByteArray)), cb->Receiver(), cb->Method());
                break;
            }
            ++iter;
        }
        m_CallbacksMux.unlock();
    }
    void Active(int cmd, QByteArray& msg){
        m_CallbacksMux.lock();
        if (!m_Callbacks.contains(cmd)){
            QList<Callback *> tmp;
            m_Callbacks[cmd] = tmp;
        }
        QList<Callback *>::iterator iter = m_Callbacks[cmd].begin();
        while(iter != m_Callbacks[cmd].end()){
            Callback * cb = *iter;
            emit cb->sigActive(msg);
            ++iter;
        }
        m_CallbacksMux.unlock();
    }

private:
    explicit CGoInterface(QObject *parent = Q_NULLPTR);
    ~CGoInterface();

private:
    static CGoInterface* m_iInstance;
    QMap<int, QList<Callback *> > m_Callbacks;
    QMutex m_CallbacksMux;
};

#endif // __GO_INTERFACE_H__

