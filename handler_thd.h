#ifndef CHANDLERTHD_H
#define CHANDLERTHD_H

#include <QObject>
#include <QThread>

class CHandlerThd : public QThread
{
    Q_OBJECT
public:
    CHandlerThd(QObject *parent = nullptr);
    ~CHandlerThd();

signals:

protected:
    void run() override;
};

#endif // CHANDLERTHD_H
