#include "models_gen_win.h"
#include "ui_models_gen_win.h"
#include <QHostAddress>
#include <QDebug>
#include "go_interface.h"
#include <QtEndian>
#include <QFileDialog>
#include <QMessageBox>
#include <QProcess>
#include "settings.h"
#include "page_win.h"

CModelsGenWin::CModelsGenWin(QWidget *parent)
    : QWidget(parent)
    , ui(new Ui::CModelsGenWin)
{
    m_strLastModelCodeSavePath = "";
    setWindowIcon(QIcon("://logo.ico"));
    ui->setupUi(this);

    QStringList strs = {tr("表名"),tr("说明")};
    ui->m_iPageWin->SetHeader(strs);
    __ClearUI();
//    ui->m_iTabGBox->setEnabled(false);
//    ui->m_iShowGBox->setEnabled(false);
    QVariantMap dbsetting = Settings::GetInstance()->value("DB").toMap();
    ui->m_iHostEdit->setText(dbsetting.value("Host").toString());
    ui->m_iPortSpinBox->setValue(dbsetting.value("Port").toInt());
    ui->m_iUsernameEdit->setText(dbsetting.value("UserName").toString());
    ui->m_iPasswdEdit->setText(dbsetting.value("Passwd").toString());
    ui->m_iDatabaseNameEdit->setText(dbsetting.value("DBName").toString());
    __OpenDB();
    GetPageFunCB fun = std::bind(&CModelsGenWin::__GetPage, this, std::placeholders::_1,std::placeholders::_2,std::placeholders::_3);
    ui->m_iPageWin->SetGetPageFunCB(fun);
    connect(ui->m_iPageWin, SIGNAL(sigDataSelected(QStringList,int)), this, SLOT(__OnTabSelected(QStringList,int)));
}

CModelsGenWin::~CModelsGenWin()
{
    delete ui;
}


void CModelsGenWin::__OnError()
{
    __ClearUI();
}

void CModelsGenWin::__OpenDB()
{
    QPushButton* pBtn = qobject_cast<QPushButton*>(sender());
    if(ui->m_iHostEdit->text() == ""){
        if (pBtn != nullptr)QMessageBox::warning(this, tr("打开连接错误"), tr("host不能为空"));
        ui->m_iHostEdit->setFocus();
        return;
    }
    if(ui->m_iPortSpinBox->value() == 0){
        if (pBtn != nullptr)QMessageBox::warning(this, tr("打开连接错误"), tr("端口号不能为零"));
        ui->m_iPortSpinBox->setFocus();
        return;
    }
    if(ui->m_iUsernameEdit->text() == ""){
        if (pBtn != nullptr)QMessageBox::warning(this, tr("打开连接错误"), tr("用户名不能为空"));
        ui->m_iUsernameEdit->setFocus();
        return;
    }
    if(ui->m_iPasswdEdit->text() == ""){
        if (pBtn != nullptr)QMessageBox::warning(this, tr("打开连接错误"), tr("密码不能为空"));
        ui->m_iPasswdEdit->setFocus();
        return;
    }
    if(ui->m_iDatabaseNameEdit->text() == ""){
        if (pBtn != nullptr)QMessageBox::warning(this, tr("打开连接错误"), tr("数据库名不能为空"));
        ui->m_iDatabaseNameEdit->setFocus();
        return;
    }
    pbapi::PK_OPEN_DB_REQ req;
    pbapi::PK_OPEN_DB_RSP rsp;
    req.set_host(ui->m_iHostEdit->text().toStdString());
    req.set_port(ui->m_iPortSpinBox->value());
    req.set_username(ui->m_iUsernameEdit->text().toStdString());
    req.set_passwd(ui->m_iPasswdEdit->text().toStdString());
    req.set_dbname(ui->m_iDatabaseNameEdit->text().toStdString());
    if (Call_Go_Func(&req, &rsp, OpenDb) == 0){
        if (rsp.errmsg() != "") {
            qCritical() << "call OpenDb return failed:" << rsp.errmsg().c_str();
            if (pBtn != nullptr)QMessageBox::warning(this, tr("打开连接错误"), "错误信息:"+QString::fromStdString(rsp.errmsg()));
            return;
        } else {
            QVariantMap dbsetting = Settings::GetInstance()->value("DB").toMap();
            dbsetting["Host"] = ui->m_iHostEdit->text();
            dbsetting["Port"] = ui->m_iPortSpinBox->value();
            dbsetting["UserName"] = ui->m_iUsernameEdit->text();
            dbsetting["Passwd"] = ui->m_iPasswdEdit->text();
            dbsetting["DBName"] = ui->m_iDatabaseNameEdit->text();
            Settings::GetInstance()->setValue("DB", dbsetting);
            ui->m_iPageWin->Reset();
        }
    } else {
        qCritical() << "call OpenDb failed";
        __OnError();
        if (pBtn != nullptr)QMessageBox::warning(this, tr("打开连接错误"), tr("call OpenDb failed"));
        return;
    }
}

void CModelsGenWin::__OnSearch()
{
    if(ui->m_iSearchEdit->text() == ""){
        QMessageBox::warning(this, tr("搜索错误"), tr("搜索字符串是必须的"));
        ui->m_iSearchEdit->setFocus();
        return;
    }
    ui->m_iPageWin->Reset();
}

void CModelsGenWin::__OnTabSelected(QStringList rowdata, int col)
{
    qDebug() << "selected history is " << rowdata;
    m_CurrentTabDatas = rowdata;
    if (m_CurrentTabDatas.size() < 1 || col < 0){
        m_CurrentTabDatas.clear();
        ui->m_iShowEdit->clear();
        return;
    }
    if (ui->m_iShowSqlRBtn->isChecked()){
        pbapi::PK_GET_TAB_SQL_REQ req;
        pbapi::PK_GET_TAB_SQL_RSP rsp;
        req.set_tabname(rowdata[0].toStdString());
        if (Call_Go_Func(&req, &rsp, GetTabSql) == 0) {
            qDebug() << "get tab sql response:" << rsp.errmsg().c_str();
            if (rsp.errmsg() == ""){
                ui->m_iShowEdit->clear();
                ui->m_iShowEdit->appendPlainText(rsp.sql().c_str());
            }
        } else {
            qCritical() << "rpc call GetTabSql failed" ;
        }
    } else if(ui->m_iShowModelRBtn->isChecked()){
        pbapi::PK_GET_TAB_MODEL_CODE_REQ req;
        pbapi::PK_GET_TAB_MODEL_CODE_RSP rsp;
        req.set_tabname(rowdata[0].toStdString());
        if (Call_Go_Func(&req, &rsp, GetTabModelCode) == 0) {
            qDebug() << "get tab model code response:" << rsp.errmsg().c_str();
            if (rsp.errmsg() == ""){
                ui->m_iShowEdit->clear();
                ui->m_iShowEdit->appendPlainText(rsp.code().c_str());
                ui->m_iExportModelCodeBtn->setVisible(true);
            }
        } else {
            qCritical() << "rpc call GetTabModelCode failed" ;
        }
    }
}

void CModelsGenWin::__OnShowRBtnClicked()
{
    ui->m_iShowEdit->clear();
    if (m_CurrentTabDatas.size() > 0){
        QRadioButton* pBtn = qobject_cast<QRadioButton*>(sender());
        if(pBtn == ui->m_iShowModelRBtn){
            pbapi::PK_GET_TAB_MODEL_CODE_REQ req;
            pbapi::PK_GET_TAB_MODEL_CODE_RSP rsp;
            req.set_tabname(m_CurrentTabDatas[0].toStdString());
            if (Call_Go_Func(&req, &rsp, GetTabModelCode) == 0) {
                qDebug() << "get tab model code response:" << rsp.errmsg().c_str();
                if (rsp.errmsg() == ""){
                    ui->m_iShowEdit->clear();
                    ui->m_iShowEdit->appendPlainText(rsp.code().c_str());
                    ui->m_iExportModelCodeBtn->setVisible(true);
                }
            } else {
                qCritical() << "rpc call GetTabModelCode failed";
            }
            ui->m_iExportModelCodeBtn->setVisible(true);
        } else if(pBtn == ui->m_iShowSqlRBtn){
            pbapi::PK_GET_TAB_SQL_REQ req;
            pbapi::PK_GET_TAB_SQL_RSP rsp;
            req.set_tabname(m_CurrentTabDatas[0].toStdString());
            if (Call_Go_Func(&req, &rsp, GetTabSql) == 0) {
                qDebug() << "get tab sql response:" << rsp.errmsg().c_str();
                if (rsp.errmsg() == ""){
                    ui->m_iShowEdit->clear();
                    ui->m_iShowEdit->appendPlainText(rsp.sql().c_str());
                }
            } else {
                qCritical() << "rpc call GetTabSql failed" ;
            }
            ui->m_iExportModelCodeBtn->setVisible(false);
        }
    }

}

void CModelsGenWin::__OnExportModelCode()
{
    qDebug() << "#####" <<m_strLastModelCodeSavePath;
    if (m_CurrentTabDatas.size() < 1){
        return;
    }
    QString tbname =m_CurrentTabDatas.at(0);
    QString tbprefix = ui->m_iTabPrefixEdit->text();
    if (tbprefix != "" && tbname.startsWith(tbprefix, Qt::CaseInsensitive)){
        tbname = tbname.right(tbname.size()-tbprefix.size());
    }
    tbname = (m_strLastModelCodeSavePath.size()>0?(m_strLastModelCodeSavePath+"/"):"") + tbname + ".go";
    QString fileName = QFileDialog::getSaveFileName(this, tr("Save File"),
                               tbname,
                               tr("go file (*.go)"));
    if (fileName != ""){
        QFileInfo info(fileName);
        qDebug() << "fileName=" <<fileName;
        m_strLastModelCodeSavePath = info.absolutePath();
        qDebug() << "m_strLastModelCodeSavePath=" <<m_strLastModelCodeSavePath;
        QFile file(fileName);
        if(!file.open(QIODevice::WriteOnly | QIODevice::Text | QIODevice::Truncate)){
            QMessageBox::warning(this, tr("打开文件错误"), tr("打开文件错误"));
            return;
        }
        file.write(ui->m_iShowEdit->toPlainText().toUtf8());
        file.close();
    }
}

void CModelsGenWin::__ClearUI()
{
    ui->m_iSearchEdit->clear();
    ui->m_iShowEdit->clear();
    if(ui->m_iShowModelRBtn->isChecked()){
        ui->m_iExportModelCodeBtn->setVisible(true);
    } else {
        ui->m_iExportModelCodeBtn->setVisible(false);
    }
}

QList<QStringList> CModelsGenWin::__GetPage(int page, int limit, int64_t *total)
{
    QList<QStringList> ret;
    if (page > 0){
        pbapi::PK_GET_TABNAMES_PAGE_REQ req;
        pbapi::PK_GET_TABNAMES_PAGE_RSP rsp;
        if (ui->m_iSearchEdit->text() != ""){
            req.set_filter(ui->m_iSearchEdit->text().toStdString());
        } else if(ui->m_iTabPrefixEdit->text() != ""){
            req.set_filter(ui->m_iTabPrefixEdit->text().toStdString());
        }
        req.set_limit(limit);
        req.set_page(page-1);
        if (Call_Go_Func(&req, &rsp, GetTabNamesPage) == 0) {
            if (rsp.errmsg() == ""){
                for (auto index = 0; index < rsp.names_size(); index++){
                    QStringList tmp;
                    tmp << QString::fromStdString(rsp.names(index));
                    tmp << QString::fromStdString(rsp.deses(index));
                    ret.append(tmp);
                }
                if (total != nullptr){
                    *total = rsp.total();
                }
            }
        } else {
            qCritical() << "rpc call GetTabNamesPage failed";
        }
    }
    return ret;
}


