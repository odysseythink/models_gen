#include "page_win.h"
#include "ui_page_win.h"
#include <QMessageBox>
#include <QPushButton>
#include <QDebug>

CPageWin::CPageWin(QWidget *parent) :
    QWidget(parent),
    ui(new Ui::CPageWin)
{
    ui->setupUi(this);
    m_GetPageFunCB = nullptr;
    ui->m_iDataView->setModel(&m_DataModel);
    m_nTotal = 0;
    m_nCurrentPage = 1;
    m_nWantedPage = 0;
    m_PageBtns.append(ui->m_iFirstPageBtn);
    m_PageBtns.append(ui->m_iPage1Btn);
    m_PageBtns.append(ui->m_iPage2Btn);
    m_PageBtns.append(ui->m_iPage3Btn);
    m_PageBtns.append(ui->m_iPage4Btn);
    m_PageBtns.append(ui->m_iPage5Btn);
    m_PageBtns.append(ui->m_iLastPageBtn);
}

CPageWin::~CPageWin()
{
    delete ui;
}

void CPageWin::SetGetPageFunCB(GetPageFunCB cb)
{
    m_GetPageFunCB = cb;
    __Refresh();
}

void CPageWin::Reset()
{
    m_nCurrentPage = 1;
    ui->m_iPage1Btn->setText(QString::number(2));
    if(m_GetPageFunCB !=nullptr)__Refresh();

}

void CPageWin::__ShowPageBtn()
{
    ui->m_iTotalLbl->setText(QString("共%1条").arg(m_nTotal));
    int nTotalPage = ui->m_iLastPageBtn->text().toInt();
    qDebug("m_nCurrentPage=%d, m_nTotalPage=%d", m_nCurrentPage, nTotalPage);
    if (nTotalPage > 1){
        ui->m_iPrePageBtn->setEnabled(true);
        ui->m_iNextPageBtn->setEnabled(true);
        ui->m_iFirstPageBtn->setVisible(true);
        ui->m_iLastPageBtn->setVisible(true);
    }else {
        ui->m_iPrePageBtn->setEnabled(false);
        ui->m_iNextPageBtn->setEnabled(false);
        ui->m_iFirstPageBtn->setVisible(true);
        ui->m_iLastPageBtn->setVisible(false);
    }
    for(int iLoop = 1; iLoop < m_PageBtns.size()-1; iLoop++){
        if (iLoop <= nTotalPage-2){
            m_PageBtns.at(iLoop)->setVisible(true);
        } else {
            m_PageBtns.at(iLoop)->setVisible(false);
        }
    }

    int nPage1 = ui->m_iPage1Btn->text().toInt();
    int nFirstPage = ui->m_iFirstPageBtn->text().toInt();
    if (!ui->m_iPage1Btn->isHidden() && ((nPage1 - nFirstPage) > 1)){
        ui->m_iLeftPlaceholderLbl->setVisible(true);
    } else {
        ui->m_iLeftPlaceholderLbl->setVisible(false);
    }

    int nPage5 = ui->m_iPage5Btn->text().toInt();
    qDebug() << "------ui->m_iPage5Btn->isHidden()=" << ui->m_iPage5Btn->isHidden();
    qDebug() << "------nTotalPage -nPage5=" << nTotalPage -nPage5;
    if (!ui->m_iPage5Btn->isHidden() && ((nTotalPage -nPage5) > 1)){
        qDebug() << "------m_iRightPlaceholderLbl show";
        ui->m_iRightPlaceholderLbl->setVisible(true);
    } else {
        qDebug() << "------m_iRightPlaceholderLbl not show";
        ui->m_iRightPlaceholderLbl->setVisible(false);
    }
    if (m_nCurrentPage <= nFirstPage){
        ui->m_iPrePageBtn->setEnabled(false);
    } else {
        ui->m_iPrePageBtn->setEnabled(true);
    }
    if (m_nCurrentPage >= nTotalPage){
        ui->m_iNextPageBtn->setEnabled(false);
    } else {
        ui->m_iNextPageBtn->setEnabled(true);
    }
    for(int iLoop = 0; iLoop < m_PageBtns.size(); iLoop++){
        int nPageVal = m_PageBtns.at(iLoop)->text().toInt();
        if (m_nCurrentPage == nPageVal){
            m_PageBtns.at(iLoop)->setStyleSheet ("color: red;");
        }else{
            m_PageBtns.at(iLoop)->setStyleSheet ("color: black;");
        }
    }
}

void CPageWin::__ShowPageBtnValue()
{
    int nTotalPage = ui->m_iLastPageBtn->text().toInt();
    int nFirstPage = ui->m_iFirstPageBtn->text().toInt();
    int nPage1 = ui->m_iPage1Btn->text().toInt();
    if (m_nCurrentPage > nPage1+4 && m_nCurrentPage < nTotalPage){
        nPage1 = m_nCurrentPage-4;
    } else if (m_nCurrentPage < nPage1 && m_nCurrentPage > nFirstPage){
        nPage1 = m_nCurrentPage;
    }
    for(int iLoop = 1; iLoop < m_PageBtns.size()-1; iLoop++){
        m_PageBtns.at(iLoop)->setText(QString::number(nPage1+iLoop-1));
    }
}

void CPageWin::__Refresh()
{
    qDebug()<< "---m_nCurrentPage=" << m_nCurrentPage;
    int limit = (ui->m_iLimitEdit->currentIndex()+1)*10;
    QList<QStringList> res = m_GetPageFunCB(m_nCurrentPage, limit, &m_nTotal);
    while(m_DataModel.rowCount() > 0){
        m_DataModel.removeRow(0);
    }
    int cCount = m_DataModel.columnCount();
    for (auto iLoop = 0; iLoop < res.size(); iLoop++){
        QList<QStandardItem*> rowItems;
        for (auto jLoop = 0; jLoop < cCount; jLoop++){
            QStandardItem* pItem = new QStandardItem(res.at(iLoop).at(jLoop));
            rowItems.append(pItem);
        }
        m_DataModel.appendRow(rowItems);
    }
//    ui->m_iDataView->resizeColumnsToContents();
    qDebug()<< "---m_nTotal=" << m_nTotal;
    qDebug()<< "---m_nTotal%limit=" << m_nTotal%limit;
    int totalPage = ((limit-m_nTotal%limit)+m_nTotal)/limit;
    qDebug()<< "---totalPage=" << totalPage;
    ui->m_iLastPageBtn->setText(QString::number(totalPage));
    ui->m_iFirstPageBtn->setText(QString::number(1));
    __ShowPageBtnValue();
    __ShowPageBtn();
    QStringList datas;
    emit sigDataSelected(datas, -1);
}

void CPageWin::__OnPageChange()
{
    int totalPage = ui->m_iLastPageBtn->text().toInt();
    QPushButton* pBtn = qobject_cast<QPushButton*>(sender());
    if (pBtn == ui->m_iPrePageBtn){
        m_nCurrentPage--;
        __Refresh();
    } else if (pBtn == ui->m_iNextPageBtn){
        m_nCurrentPage++;
        __Refresh();
    } else {
        int nPage = pBtn->text().toInt();
        m_nCurrentPage = nPage;
        __Refresh();
    }
}

void CPageWin::__OnLimitChange(int index)
{
    m_nCurrentPage = 1;
    __Refresh();
}

void CPageWin::__OnDataSelected(const QModelIndex &index)
{
    qDebug() << "selected history is " << m_DataModel.itemFromIndex(index)->text();
    qDebug() << "---------row=" << index.row();
    qDebug() << "---------col=" << index.column();
    QStringList datas;
    int cCount = m_DataModel.columnCount();
    for (auto iLoop = 0; iLoop < cCount; iLoop++){
        datas.append(m_DataModel.item(index.row(),iLoop)->text());
    }
    emit sigDataSelected(datas, index.column());
}


