/****************************************************************************
** Meta object code from reading C++ file 'CutyCapt.hpp'
**
** Created: Thu Apr 26 13:22:09 2012
**      by: The Qt Meta Object Compiler version 62 (Qt 4.6.2)
**
** WARNING! All changes made in this file will be lost!
*****************************************************************************/

#include "CutyCapt.hpp"
#if !defined(Q_MOC_OUTPUT_REVISION)
#error "The header file 'CutyCapt.hpp' doesn't include <QObject>."
#elif Q_MOC_OUTPUT_REVISION != 62
#error "This file was generated using the moc from 4.6.2. It"
#error "cannot be used with the include files from this version of Qt."
#error "(The moc has changed too much.)"
#endif

QT_BEGIN_MOC_NAMESPACE
static const uint qt_meta_data_CutyPage[] = {

 // content:
       4,       // revision
       0,       // classname
       0,    0, // classinfo
       0,    0, // methods
       0,    0, // properties
       0,    0, // enums/sets
       0,    0, // constructors
       0,       // flags
       0,       // signalCount

       0        // eod
};

static const char qt_meta_stringdata_CutyPage[] = {
    "CutyPage\0"
};

const QMetaObject CutyPage::staticMetaObject = {
    { &QWebPage::staticMetaObject, qt_meta_stringdata_CutyPage,
      qt_meta_data_CutyPage, 0 }
};

#ifdef Q_NO_DATA_RELOCATION
const QMetaObject &CutyPage::getStaticMetaObject() { return staticMetaObject; }
#endif //Q_NO_DATA_RELOCATION

const QMetaObject *CutyPage::metaObject() const
{
    return QObject::d_ptr->metaObject ? QObject::d_ptr->metaObject : &staticMetaObject;
}

void *CutyPage::qt_metacast(const char *_clname)
{
    if (!_clname) return 0;
    if (!strcmp(_clname, qt_meta_stringdata_CutyPage))
        return static_cast<void*>(const_cast< CutyPage*>(this));
    return QWebPage::qt_metacast(_clname);
}

int CutyPage::qt_metacall(QMetaObject::Call _c, int _id, void **_a)
{
    _id = QWebPage::qt_metacall(_c, _id, _a);
    if (_id < 0)
        return _id;
    return _id;
}
static const uint qt_meta_data_CutyCapt[] = {

 // content:
       4,       // revision
       0,       // classname
       0,    0, // classinfo
       6,   14, // methods
       0,    0, // properties
       0,    0, // enums/sets
       0,    0, // constructors
       0,       // flags
       0,       // signalCount

 // slots: signature, parameters, type, tag, flags
      13,   10,    9,    9, 0x08,
      36,    9,    9,    9, 0x08,
      61,    9,    9,    9, 0x08,
      93,    9,    9,    9, 0x08,
     103,    9,    9,    9, 0x08,
     126,  113,    9,    9, 0x08,

       0        // eod
};

static const char qt_meta_stringdata_CutyCapt[] = {
    "CutyCapt\0\0ok\0DocumentComplete(bool)\0"
    "InitialLayoutCompleted()\0"
    "JavaScriptWindowObjectCleared()\0"
    "Timeout()\0Delayed()\0reply,errors\0"
    "handleSslErrors(QNetworkReply*,QList<QSslError>)\0"
};

const QMetaObject CutyCapt::staticMetaObject = {
    { &QObject::staticMetaObject, qt_meta_stringdata_CutyCapt,
      qt_meta_data_CutyCapt, 0 }
};

#ifdef Q_NO_DATA_RELOCATION
const QMetaObject &CutyCapt::getStaticMetaObject() { return staticMetaObject; }
#endif //Q_NO_DATA_RELOCATION

const QMetaObject *CutyCapt::metaObject() const
{
    return QObject::d_ptr->metaObject ? QObject::d_ptr->metaObject : &staticMetaObject;
}

void *CutyCapt::qt_metacast(const char *_clname)
{
    if (!_clname) return 0;
    if (!strcmp(_clname, qt_meta_stringdata_CutyCapt))
        return static_cast<void*>(const_cast< CutyCapt*>(this));
    return QObject::qt_metacast(_clname);
}

int CutyCapt::qt_metacall(QMetaObject::Call _c, int _id, void **_a)
{
    _id = QObject::qt_metacall(_c, _id, _a);
    if (_id < 0)
        return _id;
    if (_c == QMetaObject::InvokeMetaMethod) {
        switch (_id) {
        case 0: DocumentComplete((*reinterpret_cast< bool(*)>(_a[1]))); break;
        case 1: InitialLayoutCompleted(); break;
        case 2: JavaScriptWindowObjectCleared(); break;
        case 3: Timeout(); break;
        case 4: Delayed(); break;
        case 5: handleSslErrors((*reinterpret_cast< QNetworkReply*(*)>(_a[1])),(*reinterpret_cast< QList<QSslError>(*)>(_a[2]))); break;
        default: ;
        }
        _id -= 6;
    }
    return _id;
}
QT_END_MOC_NAMESPACE
