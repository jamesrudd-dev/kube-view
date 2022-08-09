import GetImages from './components/GetImages';
import React from 'react';
import './App.css';
import NamespaceDropdown from './components/NamespaceDropdown';
import ClusterDropdown from './components/ClusterDropdown';
import ClusterRefresh from './components/ClusterRefresh';

function App() {
  return (
    <div className="App">
      <h1>Kube View</h1>
      <div class="row">
        <div class="column"><ClusterDropdown /></div>
        <div class="column"><NamespaceDropdown /></div>
        <div class="column"><ClusterRefresh /></div>
      </div>
      <br></br>
      <GetImages />
    </div>
  );
}

export default App;
