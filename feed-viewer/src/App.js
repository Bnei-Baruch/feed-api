import React from 'react';
import './App.css';

import FeedContainer from './FeedContainer.js'

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <p>
          Feed Viewer
        </p>
		<FeedContainer />
      </header>
    </div>
  );
}

export default App;
