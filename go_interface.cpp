#include "go_interface.h"
#include <QDebug>



int Call_Go_Func(google::protobuf::Message* req, google::protobuf::Message* rsp, GO_Interface_Func_Type callback) {
    if (req == nullptr || rsp == nullptr || callback == nullptr){
        qCritical() << "invalid args";
        return -1;
    }

    std::string reqraw;
    if (!req->SerializeToString(&reqraw)){
        qCritical() << "SerializeToString failed";
        return -2;
    }
    int l = 0;
    void* p = callback((void *)reqraw.c_str(), reqraw.size(), &l);
    if (!rsp->ParseFromArray(p, l)){
        qCritical() << "ParseFromArray failed";
        return -3;
    }
    free(p);
    return 0;
}





void C_Callback(int cmd,void* p,int l)
{
    QByteArray msg;
    msg.resize(l);
    for(int iLoop = 0; iLoop < l; iLoop++){
        msg[iLoop] = ((char *)p)[iLoop];
    }
    free(p);
    CGoInterface::GetInstance()->Active(cmd, msg);
}

CGoInterface* CGoInterface::m_iInstance = nullptr;

CGoInterface::CGoInterface(QObject *parent)
{

}

CGoInterface::~CGoInterface()
{
    QMap<int, QList<Callback *> >::iterator iter = m_Callbacks.begin();
    while(iter != m_Callbacks.end()){
        QList<Callback *>::iterator iter1 = m_Callbacks[iter.key()].begin();
        while(iter1 != m_Callbacks[iter.key()].end()){
            delete *iter1;
            ++iter1;
        }
        iter.value().clear();
        ++iter;
    }
    m_Callbacks.clear();
}
