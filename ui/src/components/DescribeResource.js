// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

import React, {Component} from 'react';
import { ArrowLeftIcon, DocumentIcon, FolderIcon } from '@heroicons/react/24/outline'

function classNames(...classes) {
    return classes.filter(Boolean).join(' ')
  }

class DescribeResource extends Component {
    constructor(props) {
        super(props);
    }

    row(k,v) {
        if(typeof v !== "undefined") {
            return this.printRow(k, v, null, 2); 
        }
    }

    twoColumnsRow(k,c1,c2) {
        if(typeof c1 !== "undefined" && typeof c2 !== "undefined") {
            return this.printRow(k, c1, c2); 
        }
    }

    multiRowsObject(k,v,filter = [], pre = false) {
        if(typeof v !== "undefined" && v != null && Object.keys(v).length > 0) {
            return Object.keys(v).map(key =>
                filter.length > 0
                ? filter.includes(key) && this.printRow(key == filter[0] ? k : null, key, v[key], 1, pre)
                : this.printRow(Object.keys(v)[0] == key ? k : null, key, v[key], 1, pre)
            )
        }
    }

    multiRowsArray(k,v,filter = []) {
        if(typeof v !== "undefined" && v.length > 0) {
            return v.map(elm =>
                this.multiRowsObject(k,elm,filter)
            )
        }
    }

    multiRowsPivot(k,v,filter = []) {
        if(typeof v !== "undefined" && v.length > 0) {
            return v.map((elm, idx) =>
                this.printRow(idx == 0 ? k : null, elm[filter[0]], elm[filter[1]])
            )
        }
    }

    printRow(c1, c2, c3, colspan = 1, pre = false) {
        return  <tr key={c2} className="border-t border-gray-200">
                    <td className="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-3">
                        {c1}
                    </td>
                    <td colSpan={colspan} className="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-3">
                        {c2}
                    </td>
                    {colspan == 1 &&
                    <td className="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-3">
                        {typeof c3 === "object" 
                        ? <table><tbody>{this.multiRowsObject("", c3)}</tbody></table> 
                        : pre ? this.printData(c3) : this.printValue(c3)}
                    </td>
                    }
                </tr>
    }

    printValue(val) {
        return typeof val !== "undefined" && val.toString();
    }

    printData(val) {
        return typeof val !== "undefined" && <pre>{val.toString()}</pre>;
    }

    get(t, prop) {
        if(typeof this.props.resource[t] !== "undefined" && typeof prop === "undefined") {
            return this.props.resource[t];
        }
        if(typeof this.props.resource[t] !== "undefined" && this.props.resource[t][prop] !== "undefined") {
            return this.props.resource[t][prop];
        }
    }

    resourceMeta() {
        return <>
            {this.row("Name", this.get("metadata", "name"))}
            {this.row("Namespace", this.get("metadata", "namespace"))}
            </>
    }

    resourceAttr() {
        switch(this.props.kind) {
            case "pod":
                return <>
                {this.row("Priority", this.get("spec", "priority"))}
                {this.row("Priority Class", this.get("spec", "priorityClassName"))}
                {this.row("Runtime Class", this.get("spec", "runtimeClassName"))}
                {this.row("Service Account", this.get("spec", "serviceAccountName"))}
                {this.twoColumnsRow("Node", this.get("spec", "nodeName"), this.get("status", "hostIP"))}
                {this.multiRowsObject("Labels", this.get("metadata", "labels"))}
                {this.multiRowsObject("Annotations", this.get("metadata", "annotations"))}
                {this.row("Start Time", this.get("status", "startTime"))}
                {this.row("Status", this.get("status", "phase"))}
                {this.row("Reason", this.get("status", "reason"))}
                {this.row("Message", this.get("status", "message"))}
                {this.multiRowsObject("Security Context", this.get("spec", "securityContext"))}
                {this.row("IP", this.get("status", "podIP"))}
                {this.multiRowsPivot("Controlled By", this.get("metadata", "ownerReferences"), ["kind", "name"])}
                {this.multiRowsArray("Containers", this.get("spec", "containers"), ["name", "image", "ports", "resources", "securityContext", "volumeMounts"])}
                {this.multiRowsPivot("Conditions", this.get("status", "conditions"), ["type", "status"])}
                {this.multiRowsArray("Volumes", this.get("spec", "volumes"))}
                </>
            case "svc":
                return <>
                {this.row("Cluster IP", this.get("spec", "clusterIP"))}
                {this.multiRowsPivot("Ports", this.get("spec", "ports"), ["name", "port"])}
                {this.multiRowsObject("Selector", this.get("spec", "selector"))}
                {this.multiRowsObject("Labels", this.get("metadata", "labels"))}
                {this.multiRowsObject("Annotations", this.get("metadata", "annotations"))}
                </>
            case "deploy":
                return <>
                {this.multiRowsPivot("Conditions", this.get("status", "conditions"), ["type", "status"])}
                {this.multiRowsObject("Strategy", this.get("spec", "strategy"))}
                {this.multiRowsObject("Selector", this.get("spec", "selector"))}
                {this.multiRowsObject("Labels", this.get("metadata", "labels"))}
                {this.multiRowsObject("Annotations", this.get("metadata", "annotations"))}
                </>
            case "rs":
                return <>
                {this.multiRowsPivot("Conditions", this.get("status", "conditions"), ["type", "status"])}
                {this.row("Replicas", this.get("spec", "replicas"))}
                {this.multiRowsObject("Selector", this.get("spec", "selector"))}
                {this.multiRowsObject("Labels", this.get("metadata", "labels"))}
                {this.multiRowsObject("Annotations", this.get("metadata", "annotations"))}
                </>
            case "cm":
                return <>
                {this.multiRowsPivot("Conditions", this.get("status", "conditions"), ["type", "status"])}
                {this.multiRowsObject("Data", this.get("data"), [], true)}
                {this.multiRowsObject("Labels", this.get("metadata", "labels"))}
                {this.multiRowsObject("Annotations", this.get("metadata", "annotations"))}
                </>
            case "crb":
                return <>
                {this.multiRowsObject("Role Ref", this.get("roleRef"))}
                {this.multiRowsArray("Subjects", this.get("subjects"))}
                {this.multiRowsObject("Labels", this.get("metadata", "labels"))}
                {this.multiRowsObject("Annotations", this.get("metadata", "annotations"))}
                </>
            case "cronjob":
                return <>
                {this.multiRowsObject("Status", this.get("status"))}
                {this.multiRowsObject("Spec", this.get("spec"))}
                {this.multiRowsObject("Labels", this.get("metadata", "labels"))}
                {this.multiRowsObject("Annotations", this.get("metadata", "annotations"))}
                </>
          }
    }

    render() {
        return (
            <div className="relative mt-6 flex-1 px-4 sm:px-6">
                <table className="table-auto w-300">
                    <tbody className="bg-white">
                        {this.resourceMeta()}
                        {this.resourceAttr()}
                    </tbody>
                </table>
            </div>
        );
    }
}

export default DescribeResource;