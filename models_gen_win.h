#ifndef CSQLPUBLISH_H
#define CSQLPUBLISH_H

#include <QWidget>
#include <QTcpSocket>
#include <QTimer>
#include <QList>
#include <QProcess>
#include <iostream>
#include <memory>
#include <string>




#define MAX_PKG_LEN 1024*5
#define MAX_RECV_BUFF_LEN 1024*10

QT_BEGIN_NAMESPACE
namespace Ui { class CModelsGenWin; }
QT_END_NAMESPACE

class CModelsGenWin : public QWidget
{
    Q_OBJECT

public:
    CModelsGenWin(QWidget *parent = nullptr);
    ~CModelsGenWin();

private slots:
    void __OnError();
    void __OpenDB();
    void __OnSearch();
    void __OnTabSelected(QStringList rowdata, int col);
    void __OnShowRBtnClicked();
    void __OnExportModelCode();

private:
    void __ClearUI();
    QList<QStringList> __GetPage(int page, int limit, int64_t* total);

signals:


private:
    Ui::CModelsGenWin *ui;
    QString m_strLastModelCodeSavePath;
    QStringList m_CurrentTabDatas;
};
#endif // CSQLPUBLISH_H
