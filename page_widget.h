#ifndef PAGE_WIDGET_H
#define PAGE_WIDGET_H

#include <QWidget>
#include <functional>
#include <QStandardItemModel>
#include <QPushButton>
#include <QList>
#include <QMap>

using GetPageFunCB = std::function<QList<QStringList>(int, int, int64_t*)>;
using OperationFunCB = std::function<void (QStringList)>;

namespace Ui {
class CPageWidget;
}

class CPageWidget : public QWidget
{
    Q_OBJECT

public:
    explicit CPageWidget(QWidget *parent = nullptr);
    ~CPageWidget();
    void SetHeader(QStringList header){
        m_DataModel.setHorizontalHeaderLabels(header);
        Refresh();
    }
    void SetGetPageFunCB(GetPageFunCB cb);
    void Reset();
    void Refresh();
    void AddOperation(QString title, OperationFunCB cb);

private:
    void __ShowPageBtn();
    void __ShowPageBtnValue();
    void __ChangePageBtnSize();
    QSize __GetTextSize(const QString &text);

private slots:
    void __OnPageChange();
    void __OnLimitChange(int index);
    void __OnDataSelected(const QModelIndex &index);
    void __OnOpenration();

signals:
    void sigDataSelected(QStringList rowdata, int col);

private:
    Ui::CPageWidget *ui;
    int64_t m_nTotal;
    int32_t m_nCurrentPage;
    int32_t m_nWantedPage;
    QStandardItemModel m_DataModel;
    GetPageFunCB m_GetPageFunCB;
    QList<QPushButton* > m_PageBtns;
    QMap<QString, OperationFunCB> m_Operations;
};

#endif // PAGE_WIDGET_H
