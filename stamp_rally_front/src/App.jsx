import { useState } from 'react';
import './App.css';
import Home from './Home';
import Bingo from './Bingo';
import Qr from './Qr';
import Prize from './Prize';
import Manual from './Manual';

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