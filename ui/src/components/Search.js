// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

import {Component} from 'react';
import { ArrowLeftIcon } from '@heroicons/react/24/outline';
import * as k8s from './utils/K8s';

class Search extends Component {
  constructor(props) {
    super(props);
    this.state = {
        data: new Map()
    }
  }

  componentDidMount() {
    console.log("SearchView componentDidMount", this.props);
    this.search();
  }

  search() {
    for (let i=0; i < k8s.resources.length; i++) {
      this.updateList(k8s.resources[i].id);
    }
    for (let i=0; i < k8s.cluster.length; i++) {
      this.updateList(k8s.cluster[i].id);
    }
  }

  updateList(key) {
    fetch('/api/v1/list/' + key + '?filter=' + this.props.filter)
    .then(res => res.json())
    .then(values => {
      var newMap = this.state.data;
      if (values && values.length > 0) {
        newMap.set(key, values);
        this.setState((state, props) => ({
          data: newMap
        }));
      } else {
        newMap.delete(key);
        this.setState((state, props) => ({
          data: newMap
        }));
      }
    });
  }

  componentDidUpdate(prevProps, prevState) {
    if(prevState !== this.state) {
      console.log('SearchView componentDidUpdate state: ', this.state);
    }
    if(prevProps !== this.props) {
      console.log('SearchView componentDidUpdate props: ', this.props);
      this.search();
    }
  }

  render() {
    return (
          <div className="flex h-full flex-col overflow-y-scroll bg-white bg-opacity-80 py-6 shadow-xl">
            <div className="px-4 sm:px-6">
              <div className="flex items-start justify-between">
                <div className="">
                  <div className="font-bold text-xl leading-6 text-gray-900 border-b-4 p-4">
                    <button
                      type="button"
                      className="pr-4 rounded-md bg-white text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
                      onClick={this.props.close}
                    >
                      <span className="sr-only">Back</span>
                      <ArrowLeftIcon className="h-6 w-6" aria-hidden="true" />
                    </button>
                    Search
                  </div>
                </div>
              </div>
            </div>
            <div className="relative mt-6 flex-1 px-4 sm:px-6">
            {Array.from(this.state.data.keys()).length > 0 ?
              <table className="table-auto w-300">
                <thead className="bg-white">
                  <tr>
                    <th scope="col" className="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-900 sm:pl-3">
                      Name
                    </th>
                    <th scope="col" className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
                      Created at
                    </th>
                  </tr>
                </thead>
                  <tbody className="bg-white">
                  {Array.from(this.state.data.keys()).map((k) => (
                    Array.from(this.state.data.get(k)).length > 0 && this.state.data.get(k).map((e) => (
                      <tr key={e.uid} className="border-t border-gray-200">
                      <td className="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-3">
                      {e.name}
                      </td>
                      <td className="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-3">
                      {e.creationTimestamp}
                      </td>     
                    </tr>
                    ))
                  ))}
                  </tbody>
              </table>
                : ""}
            </div> 
          </div>
    );
  }
}

export default Search;