package com.scanzero.app

import com.facebook.react.bridge.ReactApplicationContext
import com.facebook.react.bridge.ReactContextBaseJavaModule
import com.facebook.react.bridge.ReactMethod
import com.facebook.react.bridge.Promise

import com.amirmiir.scannercore.scannercore.Scannercore


class ScannerCoreModule(reactContext:ReactApplicationContext): ReactContextBaseJavaModule(reactContext) {
    override fun getName():String{
        return "ScannerCore"
    }

    @ReactMethod
    fun ping(input:String, promise: Promise){
        //1.try to call Scannercore.ping(input)
        //2. resolve the promise with result
        //catch exceptions:
        //reject promise with error message
        try{
            promise.resolve(Scannercore.ping(input))
        } catch(e: Exception){
            promise.reject("E_SCANNERCORE", e)
        }
    }
}