QT       += core gui network

greaterThan(QT_MAJOR_VERSION, 4): QT += widgets

CONFIG += c++11

# You can make your code fail to compile if it uses deprecated APIs.
# In order to do so, uncomment the following line.
#DEFINES += QT_DISABLE_DEPRECATED_BEFORE=0x060000    # disables all the APIs deprecated before Qt 6.0.0

SOURCES += \
    go_interface.cpp \
    handler_thd.cpp \
    main.cpp \
    settings.cpp \
    api.pb.cc \
    models_gen_win.cpp \
    page_win.cpp

HEADERS += \
    go_interface.h \
    go_models_gen/c_2_go_interface.h \
    go_models_gen/libmodels_gen.h \
    settings.h \
    api.pb.h \
    handler_thd.h \
    models_gen_win.h \
    page_win.h

FORMS += \
    models_gen_win.ui \
    page_win.ui

# Default rules for deployment.
qnx: target.path = /tmp/$${TARGET}/bin
else: unix:!android: target.path = /opt/$${TARGET}/bin
!isEmpty(target.path): INSTALLS += target

RESOURCES += \
    res.qrc

RC_ICONS = logo.ico

unix: LIBS += "../go_models_gen/libmodels_gen.a" -ldl -lresolv
else: LIBS += "../go_models_gen/libmodels_gen.a"

LIBS += -static -lprotobuf

