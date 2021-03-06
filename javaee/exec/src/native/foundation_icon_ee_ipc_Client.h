/* DO NOT EDIT THIS FILE - it is machine generated */
#include <jni.h>
/* Header for class foundation_icon_ee_ipc_Client */

#ifndef _Included_foundation_icon_ee_ipc_Client
#define _Included_foundation_icon_ee_ipc_Client
#ifdef __cplusplus
extern "C" {
#endif
/*
 * Class:     foundation_icon_ee_ipc_Client
 * Method:    socket
 * Signature: ()I
 */
JNIEXPORT jint JNICALL Java_foundation_icon_ee_ipc_Client_socket
  (JNIEnv *, jclass);

/*
 * Class:     foundation_icon_ee_ipc_Client
 * Method:    connect
 * Signature: (ILjava/lang/String;)V
 */
JNIEXPORT void JNICALL Java_foundation_icon_ee_ipc_Client_connect
  (JNIEnv *, jclass, jint, jstring);

/*
 * Class:     foundation_icon_ee_ipc_Client
 * Method:    close
 * Signature: (I)V
 */
JNIEXPORT void JNICALL Java_foundation_icon_ee_ipc_Client_close
  (JNIEnv *, jclass, jint);

/*
 * Class:     foundation_icon_ee_ipc_Client
 * Method:    read
 * Signature: (I[BII)I
 */
JNIEXPORT jint JNICALL Java_foundation_icon_ee_ipc_Client_read
  (JNIEnv *, jclass, jint, jbyteArray, jint, jint);

/*
 * Class:     foundation_icon_ee_ipc_Client
 * Method:    write
 * Signature: (I[BII)V
 */
JNIEXPORT void JNICALL Java_foundation_icon_ee_ipc_Client_write
  (JNIEnv *, jclass, jint, jbyteArray, jint, jint);

#ifdef __cplusplus
}
#endif
#endif
