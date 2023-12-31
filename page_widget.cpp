#include "page_widget.h"
#include "ui_page_widget.h"
#include <QMessageBox>
#include <QPushButton>
#include <QDebug>
#include <QButtonGroup>
#include <QHBoxLayout>


CPageWidget::CPageWidget(QWidget *parent) :
    QWidget(parent),
    ui(new Ui::CPageWidget)
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
    __ChangePageBtnSize();
}

CPageWidget::~CPageWidget()
{
    delete ui;
}

void CPageWidget::SetGetPageFunCB(GetPageFunCB cb)
{
    m_GetPageFunCB = cb;
    Refresh();
}

void CPageWidget::Reset()
{
    m_nCurrentPage = 1;
    ui->m_iPage1Btn->setText(QString::number(2));
    if(m_GetPageFunCB !=nullptr)Refresh();

}

void CPageWidget::__ShowPageBtn()
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

void CPageWidget::__ShowPageBtnValue()
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

void CPageWidget::__ChangePageBtnSize()
{
    for(int iLoop = 0; iLoop < m_PageBtns.size(); iLoop++){
        QSize sz = __GetTextSize(m_PageBtns.at(iLoop)->text());
        if(m_PageBtns.at(iLoop) == ui->m_iLastPageBtn){
            qDebug() << "---------m_iLastPageBtn width=" << sz.width();
        }
        m_PageBtns.at(iLoop)->setMinimumWidth(sz.width()+10);
        m_PageBtns.at(iLoop)->setMaximumWidth(sz.width()+10);
    }
}

QSize CPageWidget::__GetTextSize(const QString &text)
{
//        /* 设置字体属性 */
//        QFont font;
//        font.setPixelSize(45);
//        font.setFamily("Microsoft YaHei UI");

        /* 设置字体信息 */
        QFontMetrics metrics(ui->m_iLastPageBtn->font());

        return metrics.size(Qt::TextSingleLine, text);
}

void CPageWidget::Refresh()
{
    if (m_GetPageFunCB == nullptr){
        return;
    }
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
            QStringList tmplist = res.at(iLoop);
            QStandardItem* pItem = new QStandardItem(jLoop>=tmplist.size()?"":res.at(iLoop).at(jLoop));
            rowItems.append(pItem);
        }
        m_DataModel.appendRow(rowItems);
        if (m_Operations.size() > 0){
            QWidget* pGroup = new QWidget(this);
            QHBoxLayout* pGroupLayout = new QHBoxLayout();
            pGroupLayout->setContentsMargins(0,0,0,0);
            pGroupLayout->setSpacing(5);
            pGroup->setLayout(pGroupLayout);
            pGroupLayout->addSpacerItem(new QSpacerItem(16777215, 0));
            QStringList operationNames = m_Operations.keys();
            for(auto jLoop = 0; jLoop < operationNames.size(); jLoop++){
                QPushButton* pBtn = new QPushButton(operationNames.at(jLoop));
                qDebug() << "------++" << operationNames.at(jLoop).size();
                pBtn->setFixedWidth(46);
                pBtn->setProperty("row", iLoop);
                connect(pBtn, SIGNAL(clicked()), this, SLOT(__OnOpenration()));
                pGroupLayout->addWidget(pBtn);
            }
            pGroupLayout->addSpacerItem(new QSpacerItem(16777215, 0));
            ui->m_iDataView->setIndexWidget(m_DataModel.index(iLoop, cCount-1), pGroup);
        }
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
    __ChangePageBtnSize();
    QStringList datas;
    emit sigDataSelected(datas, -1);
}

void CPageWidget::AddOperation(QString title, OperationFunCB cb)
{
    int colCount = m_DataModel.columnCount();
    QStandardItem* pItem = m_DataModel.horizontalHeaderItem(colCount-1);
    if (pItem != nullptr){
        qDebug() << "--------" << pItem->text();
        if(pItem->text() != tr("操作")){
            QStringList header;
            for (int iLoop = 0; iLoop < colCount; iLoop++){
                header << m_DataModel.horizontalHeaderItem(iLoop)->text();
            }
            header << tr("操作");
            SetHeader(header);
            m_Operations[title] = cb;
        }
    }
}

void CPageWidget::__OnPageChange()
{
    int totalPage = ui->m_iLastPageBtn->text().toInt();
    QPushButton* pBtn = qobject_cast<QPushButton*>(sender());
    if (pBtn == ui->m_iPrePageBtn){
        m_nCurrentPage--;
        Refresh();
    } else if (pBtn == ui->m_iNextPageBtn){
        m_nCurrentPage++;
        Refresh();
    } else {
        int nPage = pBtn->text().toInt();
        m_nCurrentPage = nPage;
        Refresh();
    }
}

void CPageWidget::__OnLimitChange(int index)
{
    m_nCurrentPage = 1;
    Refresh();
}

void CPageWidget::__OnDataSelected(const QModelIndex &index)
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

void CPageWidget::__OnOpenration()
{
    QPushButton* pBtn = qobject_cast<QPushButton*>(sender());
    if (pBtn != nullptr){
        if(m_Operations.contains(pBtn->text()) && m_Operations[pBtn->text()] != nullptr){
            QStringList data;
            int row = pBtn->property("row").toInt();
            int cCount = m_DataModel.columnCount();
            for (auto iLoop = 0; iLoop < cCount-1; iLoop++){
                data.append(m_DataModel.item(row,iLoop)->text());
            }
            m_Operations[pBtn->text()](data);
        }
    }
}


