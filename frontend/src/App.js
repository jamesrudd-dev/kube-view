import React from 'react';
import 'bootstrap/dist/css/bootstrap.min.css';
import './App.css';
import ClusterDropdown from './components/ClusterDropdown';

function App() {
  return (
    <div className="App">
      <h1>Kube View</h1>
      <ClusterDropdown />
    </div>
  );
}

export default App;
