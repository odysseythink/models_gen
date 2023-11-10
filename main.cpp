#include "models_gen_win.h"

#include <QApplication>
#include "settings.h"

int main(int argc, char *argv[])
{
    QApplication a(argc, argv);
    Settings::GetInstance();
    qDebug() << Settings::APP_NAME << " starting......";
    CModelsGenWin w;
    w.show();
    int nret = a.exec();
    delete Settings::GetInstance();
    return nret;
}
