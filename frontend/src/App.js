import GetImages from './components/GetImages';
import React from 'react';
import './App.css';
import NamespaceDropdown from './components/NamespaceDropdown';

function App() {
  return (
    <div className="App">
      <h1>Kube View</h1>
      <NamespaceDropdown />
      <br></br>
      <GetImages />
    </div>
  );
}

export default App;
