import { useState } from 'react';
import './App.css';
import Home from "./components/Home";
import Bingo from "./components/Bingo";
import Qr from "./components/Qr";
import Prize from "./components/Prize";
import Manual from "./components/Manual";

function App() {
  const [currentScreen, setCurrentScreen] = useState('home');

  // 画面遷移用の関数
  const navigate = (screenName) => {
    setCurrentScreen(screenName);
  };

  return (
    <div className="phone">
      {currentScreen === 'home' && <Home navigate={navigate} />}
      {currentScreen === 'bingo' && <Bingo navigate={navigate} currentScreen={currentScreen} />}
      {currentScreen === 'qr' && <Qr navigate={navigate} currentScreen={currentScreen} />}
      {currentScreen === 'prize' && <Prize navigate={navigate} currentScreen={currentScreen} />}
      {currentScreen === 'manual' && <Manual navigate={navigate} />}
    </div>
  );
}

export default App;