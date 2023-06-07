// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

import { useState, useEffect } from 'react'
import { Routes, Route } from 'react-router-dom';
import Login from './components/Login';
import Stage from './components/Stage';

function App() {
   const [isLoggedIn, setIsLoggedIn] = useState(false);
   const [isLoaded, setIsLoaded] = useState(false);
   
   useEffect(() => {
      fetch('/api/v1/auth/status', {
         method: 'GET',
      })
     .then(res => res)
      .then(d => {
         if (d.status == 200) {
            setIsLoggedIn(true);
            console.log('/api/v1/auth/status', d.status);
         }
         setIsLoaded(true);
      });
   }, []);

   return (
      <>{ isLoaded ? 
         <Routes>
            <Route path="/login" element={<Login />} />
            <Route path="/" element={!isLoggedIn ? <Login /> : <Stage />} />
            <Route path="/cluster/:cluster" element={!isLoggedIn ? <Login /> : <Stage />} />
            <Route path="/clusterresource/:kind/:resource" element={!isLoggedIn ? <Login /> : <Stage />} />
            <Route path="/node/:node" element={!isLoggedIn ? <Login /> : <Stage />} />
            <Route path="/namespace/:namespace/:kind" element={!isLoggedIn ? <Login /> : <Stage />} />
            <Route path="/resource/:namespace/:kind/:resource" element={!isLoggedIn ? <Login /> : <Stage />} />  
         </Routes>
         : null}
      </>   
   );
};

export default App;