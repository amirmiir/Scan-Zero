package com.scanzero.app

import com.facebook.react.ReactPackage
import com.facebook.react.bridge.NativeModule
import com.facebook.react.bridge.ReactApplicationContext
import com.facebook.react.uimanager.ViewManager

class ScannerCorePackage: ReactPackage{
    override fun createNativeModules(reactContext: ReactApplicationContext): List<NativeModule>{
        return listOf(ScannerCoreModule(reactContext))
    }

    override fun createViewManagers(reactContext: ReactApplicationContext): List<ViewManager<*, *>>{
        return emptyList()
    }
}