// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

import React, {Component} from 'react';
import { ArrowLeftIcon, DocumentIcon, FolderIcon } from '@heroicons/react/24/outline'

function classNames(...classes) {
    return classes.filter(Boolean).join(' ')
  }

class FilesystemBrowser extends Component {
    constructor(props) {
        super(props);
        this.state = {
            current: "",
            inodes: [],
            content: null,
        }
    }

    componentDidMount() {
      console.log("FilesystemBrowser componentDidMount");
      this.fetchUrl("");
    }

    onClick = (path, type = "d") => {
        console.log("FilesystemBrowser onClick: ", path);
        if (type === "r" || type === "s") {
          this.fetchFile(path);
        } else {
          this.fetchUrl(path);
        }
    };

    fetchFile(path) {
      const url = '/api/v1/pod/file/' + this.props.url + path;
      fetch(url)
      .then(res => res)
      .then(res => {
        if (res.status == 200) {
        res.text().then(res => {
          this.setState((state, props) => ({
            content: res,
            current: path
          }));
        });
        } else if (res.status == 500) {
          console.log("can not open file");
        }
      }).catch(function(error) {
        console.log(error);
      });
    }

    fetchUrl(path) {
      var url;
      if (path != "") {
        url = '/api/v1/pod/fs/' + this.props.url + path;
      } else {
        url = '/api/v1/pod/fs/' + this.props.url;
      }
      fetch(url)
      .then(res => res)
      .then(res => {
        if (res.status == 200) {
        res.json().then(res => {
          const sorted = res.sort((x,y) => {if (x.name < y.name) {
            return -1;
          }});
          this.setState((state, props) => ({
            inodes: sorted,
            current: path,
            content: null
          }));
        });
        } else if (res.status == 500) {
          console.log("can not open fs inode");
        }
      }).catch(function(error) {
        console.log(error);
      });
    }

    stripPath(name) {
      console.log("current: ", this.state.current);
      if (this.state.current !== "" && this.state.current !== name) {
        return name.replace(this.state.current, "");
      }
      return name;
    }

    isCurrentFolder(name) {
      if (this.state.current === name) {
        return true;
      }
      return false;
    }

    parentFolder() {
      var splitted = this.state.current.split("/");
      var r = splitted.slice(0, splitted.length-1);
      console.log("parent: ", r);
      if (r.length > 0) {
        return r.join("/");
      }
      return "/";
    }

    isRoot() {
      if (this.state.current === "") {
        return true;
      }
      return false;
    }

    humanRead(bytes) {
      const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB']
      if (bytes === 0) return 'n/a'
      const i = parseInt(Math.floor(Math.log(bytes) / Math.log(1024)), 10)
      if (i === 0) return `${bytes} ${sizes[i]}`
      return `${(bytes / (1024 ** i)).toFixed(1)} ${sizes[i]}`
    }
    
    render() {
        return (
            <div className='overflow-x-auto'>
            <div className="relative mt-6 flex-1 px-4 sm:px-6">
            {this.state.content == null ?
              <ul>
              {!this.isRoot()
              ? <li key="parent" className='flex cursor-pointer select-none items-center rounded-md px-3 py-2' onClick={() => this.onClick(this.parentFolder())}>
                  <ArrowLeftIcon
                    className={classNames('h-6 w-6 flex-none', 'text-gray-400')}
                    aria-hidden="true"
                  />
                  <span className='ml-3 flex-auto truncate'></span>
                </li>
              : ""
              } 
              {this.state.inodes && this.state.inodes.length > 0 && this.state.inodes.map((inode) => (
                <li key={inode.name} className='flex cursor-pointer select-none items-center rounded-md px-3 py-2 hover:bg-slate-200' onClick={() => this.onClick(inode.name, inode.type)}>
                {inode.type === "d"
                ? inode.name != this.state.current
                  ? <FolderIcon
                      className={classNames('h-6 w-6 flex-none', 'text-gray-400')}
                      aria-hidden="true"
                    />
                  : ""
                : <DocumentIcon
                    className={classNames('h-6 w-6 flex-none', 'text-gray-400')}
                    aria-hidden="true"
                  />
                }
                <span className={classNames('ml-3 flex-auto truncate', this.isCurrentFolder(inode.name) ? 'font-bold text-xl' : null)}>{this.stripPath(inode.name)}</span>
                <span className="w-1/6">{inode.permission}</span> 
                <span className="w-1/6">{inode.user}:{inode.group}</span>
                <span className="w-1/6">{inode.modified}</span>
                <span className="w-1/6">{inode.type === "r" && inode.size != 0 ? this.humanRead(inode.size) : null}</span>  
              </li>                
              ))}
              </ul>
              :
              <div>
              <li key="parent" className='flex cursor-pointer select-none items-center rounded-md px-3 py-2' onClick={() => this.onClick(this.parentFolder())}>
                  <ArrowLeftIcon
                    className={classNames('h-6 w-6 flex-none', 'text-gray-400')}
                    aria-hidden="true"
                  />
                  <span className='ml-3 flex-auto truncate'></span>
                </li>
              
              <span className="w-1/6"><pre>{this.state.content}</pre></span>
              </div>
              }
            </div>
          </div>
        );
    }
}

export default FilesystemBrowser;