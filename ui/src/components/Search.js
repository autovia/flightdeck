// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

import {Component} from 'react';
import { ArrowLeftIcon } from '@heroicons/react/24/outline';
import ResourceOverlay from './Resource';
import * as k8s from './utils/K8s';

class Search extends Component {
  constructor(props) {
    super(props);
    this.state = {
        data: new Map(),
        overlay: {view: false},
        resource: {}
    }
  }

  componentDidMount() {
    console.log("SearchView componentDidMount", this.props);
    this.search();
  }

  closeResourceOverlay = (e) => {
    this.setState((state, props) => ({
      overlay: {view: false}
    }));
  }

  openResourceOverlay = (e, i, k) => {
    console.log(e.target.source);
    this.setState((state, props) => ({
      overlay: {view: true},
      resource: {kind: k, label: i.name, namespace: i.namespace}
    }));
  }

  search() {
    for (let i=0; i < k8s.resources.length; i++) {
      this.updateList(k8s.resources[i].id);
    }
    for (let i=0; i < k8s.cluster.length; i++) {
      this.updateList(k8s.cluster[i].id);
    }
  }

  rowTitle(key) {
    const all = [...k8s.resources, ...k8s.cluster];
    return all.filter((f) => f.id === key)[0].name + " " + this.state.data.get(key).length;
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

  rows(k) {
    return (
      Array.from(this.state.data.get(k)).length > 0 && this.state.data.get(k).map((e) => (
        <tr key={e.uid} onClick={(event) => this.openResourceOverlay(event, e, k)} className="border-t border-gray-200">
        <td className="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-3">
        {e.name}
        </td>
        <td className="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-3">
        {e.creationTimestamp}
        </td>     
      </tr>
      ))
    )
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
            <div className='mt-6 px-4 sm:px-6'>
            {Array.from(this.state.data.keys()).map((k) => (
              <a key={"badge-" + k} href={"#"+ k} className="inline-flex items-center rounded-md bg-gray-100 m-2 px-2 py-1 text-xs font-medium text-gray-600">
                {this.rowTitle(k)}
              </a>
            ))}
            </div>

            <div className='columns-2'>
              <div className="relative mt-6 flex-1 px-4 sm:px-6 overflow-x-auto">
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
                      <>
                      <tr id={k} key={k} className="border-t border-gray-200">
                        <th colSpan={2} scope="colgroup" className="bg-gray-50 py-2 pl-4 pr-3 text-left text-sm font-semibold text-gray-900 sm:pl-3">
                        {this.rowTitle(k)}
                        </th>     
                      </tr>
                      {this.rows(k)}
                      </> 
                    ))}
                    </tbody>
                </table>
                  : ""}
              </div> 
              {this.state.overlay.view ? <ResourceOverlay overlay={false} data={this.state.resource} close={this.closeResourceOverlay} /> : <div></div>}
            </div>
          </div>
    );
  }
}

export default Search;