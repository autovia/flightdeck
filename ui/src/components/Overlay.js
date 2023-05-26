// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

import React, {Fragment, Component} from 'react';
import { Dialog, Transition } from '@headlessui/react'
import { XMarkIcon } from '@heroicons/react/24/outline'
import Tabs from './Tabs';
import FilesystemBrowser from './FilesystemBrowser';

//export default function Overlay({ data, params, close }) {
class Overlay extends Component {
  constructor(props) {
    super(props);
    this.state = {
        open: true,
        resource: null,
        references: [],
        tabs: this.initTabs(),
        currentTab: "yaml",
        currentContainer: this.props.data.kind === "pod" && this.props.data.containers && this.props.data.containers.length > 0 ? this.props.data.containers[0] : ""
    }
  }

  componentDidMount() {
    this.fetchYaml();
  }

  initTabs() {
    var tabs = [{ name: 'YAML', id: "yaml" }]
    if (this.props.data.kind === "pod") {
      tabs.push({ name: 'Logs', id: "logs" })
      tabs.push({ name: 'Files', id: "files" }) 
    } else {
      tabs.push({ name: 'Used by pods', id: "ref" }) 
    }
    return tabs;
  }

  fetchYaml() {
    this.setState((state, props) => ({
      currentTab: "yaml"
    }));
    if (this.props.data.kind === "vol") {
      this.fetchUrl('/api/v1/' + this.props.data.kind + '/'  + this.props.params.namespace + '/' + this.props.params.pod + '/'  + this.props.data.label);
    } else {
      this.fetchUrl('/api/v1/' + this.props.data.kind + '/'  + this.props.params.namespace + '/' + this.props.data.label);
    }
  }

  fetchLogs() {
    this.setState((state, props) => ({
      currentTab: "logs"
    }));
    if (this.state.currentContainer != "") {
      this.fetchUrl('/api/v1/logs/'  + this.props.params.namespace + '/' + this.props.data.label + '/' + this.state.currentContainer);
    } else {
      this.fetchUrl('/api/v1/logs/'  + this.props.params.namespace + '/' + this.props.data.label);
    }
  }

  fetchFiles() {
    this.setState((state, props) => ({
      currentTab: "files"
    }));
  }

  fetchRefs() {
    const url = '/api/v1/graph/' + this.props.data.kind + '/'  + this.props.params.namespace + '/' + this.props.data.label;
    fetch(url)
    //.then(res => res.text())
    .then(res => res.json())
    .then(d => {
      if (d.nodes && d.nodes.length > 0) {
        const pods = d.nodes.filter((f) => f.data.kind === "pod").map(m => m.data.label);
        this.setState((state, props) => ({
          references: pods,
          currentTab: "ref"
        }));
      }
    });
  }

  fetchUrl(url) {
    fetch(url)
    .then(res => res.text())
    //.then(res => res.json())
    .then(d => {
      this.setState((state, props) => ({
        resource: d
      }));
    });
  }

  onTabClick = (id) => {
    console.log("onTabClick: ", id);
    switch(id) {
      case "logs":
        this.fetchLogs();
        break;
      case "files":
          this.fetchFiles();
          break;
      case "yaml":
        this.fetchYaml();
        break;
      case "ref":
        this.fetchRefs();
        break;
    }
  }

  onContainerChange = (e) => {
    this.setState((state, props) => ({
      currentContainer: e.target.value
    }));
  }

  componentDidUpdate(prevProps, prevState) {
    if (prevState.currentContainer !== this.state.currentContainer) {
      this.fetchLogs();
    }
  }

  render() {
    return (
    <Transition.Root show={this.state.open} as={Fragment}>
      <Dialog as="div" className="relative z-10" onClose={this.props.close}>
      <div className="fixed inset-0" />
        <div className="fixed inset-0 overflow-hidden">
          <div className="absolute inset-0 overflow-hidden">
            <div className="pointer-events-none fixed bottom-0 h-[80%] flex max-w-full ">
              <Transition.Child
                as={Fragment}
                enter="transform transition ease-in-out duration-500 sm:duration-700"
                enterFrom="translate-x-full"
                enterTo="translate-x-0"
                leave="transform transition ease-in-out duration-500 sm:duration-700"
                leaveFrom="translate-x-0"
                leaveTo="translate-x-full"
              >
                <Dialog.Panel className="pointer-events-auto w-screen">
                  <div className="flex h-full flex-col overflow-y-scroll bg-white bg-opacity-80 py-6 shadow-xl">
                    <div className="px-4 sm:px-6">
                      <div className="flex items-start justify-between">
                        
                        <Dialog.Title className="">
                          <div className="">
                            <img src={'/static/k8s/' + this.props.data.kind + '-256.png'} className="w-16" />
                          </div>
                          <div className="font-bold text-xl leading-6 text-gray-900 border-b-4 p-4">
                            {this.props.data.label}
                          </div>
                          <Tabs tabs={this.state.tabs} current={this.state.currentTab} onTabClick={this.onTabClick}></Tabs>
                          {this.state.currentTab === "logs" && this.props.data.containers && this.props.data.containers.length > 1
                            ? <div>Container&nbsp;
                              <select value={this.state.currentContainer} onChange={this.onContainerChange}>
                                {this.props.data.containers.map((c) => (
                                  <option key={c}>{c}</option>
                                ))}
                              </select>
                              </div>
                            : ""
                          }
                        </Dialog.Title>
                        <div className="ml-3 flex h-7 items-center">
                          <button
                            type="button"
                            className="rounded-md bg-white text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
                            onClick={this.props.close}
                          >
                            <span className="sr-only">Close panel</span>
                            <XMarkIcon className="h-6 w-6" aria-hidden="true" />
                          </button>
                        </div>
                      </div>
                    </div>
                    {this.state.currentTab === "yaml" || this.state.currentTab === "logs"
                      ? <div className="relative mt-6 flex-1 px-4 sm:px-6"><pre>{this.state.resource}</pre></div>
                      : ""
                    }
                    {this.state.currentTab === "files"
                      ? <FilesystemBrowser url={this.props.params.namespace + '/' + this.props.data.label + '/' + this.state.currentContainer}></FilesystemBrowser>
                      : ""
                    }
                    {this.state.currentTab === "ref"
                      ? <div className="relative mt-6 flex-1 px-4 sm:px-6">
                          {this.state.references && this.state.references.length > 0 && this.state.references.map((ref) => (
                            <p key={ref}>{ref}</p>
                          ))}
                        </div>
                      : ""
                    }                    
                  </div>
                </Dialog.Panel>
              </Transition.Child>
            </div>
          </div>
        </div>
      </Dialog>
    </Transition.Root>
    );
  }
}

export default Overlay;