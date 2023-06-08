// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

import React, { useRef } from "react";

export default function Login() {
  const tokenInput = useRef(null);

  function handleClick() {
    fetch('/api/v1/auth/login', {
      method: 'POST',
      headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + tokenInput.current.value
      }
    })
    .then(res => res.json())
    .then(d => {
      window.open("/", "_self");
    });
  };

  function handleKeyDown(event) {
    if (event.key === 'Enter') {
      event.preventDefault();
      handleClick();
    }
  }

  return (
      <div className="flex min-h-full flex-1 flex-col justify-center py-12 sm:px-6 lg:px-8">
        <div className="sm:mx-auto sm:w-full sm:max-w-md">
          <h2 className="mt-6 text-center text-2xl font-bold leading-9 tracking-tight text-gray-900">
            Flightdeck
          </h2>
          <h2 className="mt-6 text-center text-2xl leading-9 tracking-tight text-gray-900">
            Sign in with token
          </h2>
        </div>

        <div className="mt-10 sm:mx-auto sm:w-full sm:max-w-[480px]">
          <div className="bg-white px-6 py-12 shadow sm:rounded-lg sm:px-12">
            <div className="space-y-6">
              <div>
                <label htmlFor="token" className="block text-sm font-medium leading-6 text-gray-900">
                  Bearer token
                </label>
                <div className="mt-2">
                  <textarea
                    id="token"
                    name="token"
                    rows="5"
                    ref={tokenInput}
                    onKeyDown={(e) => handleKeyDown(e)}
                    required
                    autoFocus="true"
                    className="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                  />
                </div>
              </div>
              <div>
                <button
                  type="submit"
                  onClick={handleClick}
                  className="flex w-full justify-center rounded-md bg-indigo-600 px-3 py-1.5 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
                >
                  Sign in
                </button>
              </div>
            </div>
        </div>
      </div>
    </div>
  )
}