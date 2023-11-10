#ifndef PAGE_WIN_H
#define PAGE_WIN_H

#include <QWidget>
#include <functional>
#include <QStandardItemModel>
#include <QPushButton>
#include <QList>

using GetPageFunCB = std::function<QList<QStringList>(int, int, int64_t*)>;

namespace Ui {
class CPageWin;
}

class CPageWin : public QWidget
{
    Q_OBJECT

public:
    explicit CPageWin(QWidget *parent = nullptr);
    ~CPageWin();
    void SetHeader(QStringList header){
        m_DataModel.setHorizontalHeaderLabels(header);
    }
    void SetGetPageFunCB(GetPageFunCB cb);
    void Reset();

private:
    void __ShowPageBtn();
    void __ShowPageBtnValue();
    void __Refresh();

private slots:
    void __OnPageChange();
    void __OnLimitChange(int index);
    void __OnDataSelected(const QModelIndex &index);

signals:
    void sigDataSelected(QStringList rowdata, int col);

private:
    Ui::CPageWin *ui;
    int64_t m_nTotal;
    int32_t m_nCurrentPage;
    int32_t m_nWantedPage;
    QStandardItemModel m_DataModel;
    GetPageFunCB m_GetPageFunCB;
    QList<QPushButton* > m_PageBtns;
};

#endif // PAGE_WIN_H
