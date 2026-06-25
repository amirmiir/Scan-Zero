/**
 * Sample React Native App
 * https://github.com/facebook/react-native
 *
 * @format
 */

import { NewAppScreen } from '@react-native/new-app-screen';
import { useEffect, useState } from 'react';
import { NativeModules, StatusBar, StyleSheet, useColorScheme, View } from 'react-native';
import {
  SafeAreaProvider,
  useSafeAreaInsets,
} from 'react-native-safe-area-context';

const { ScannerCore } = NativeModules

function App() {
  const isDarkMode = useColorScheme() === 'dark';

  const [message, setMessage] = useState("Waiting for Go core...")

  useEffect(() => {
    async function verifyBridge() {
      try {
        //call the ping
        const result = await ScannerCore.ping("Loaded")
        console.log("bridge:", result)
      } catch (error) {
        setMessage("Bridge failed: " + error)
      }
    }

    verifyBridge()
  }, [])

  return (
    <SafeAreaProvider>
      <StatusBar barStyle={isDarkMode ? 'light-content' : 'dark-content'} />
      <AppContent />
    </SafeAreaProvider>
  );
}

function AppContent() {
  const safeAreaInsets = useSafeAreaInsets();

  return (
    <View style={styles.container}>
      <NewAppScreen
        templateFileName="App.tsx"
        safeAreaInsets={safeAreaInsets}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
});

export default App;
