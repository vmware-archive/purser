import { NodeDetails } from "./node-details";

export class CloudDetails{
    cloud : string;
    cpu : number;
    memory : number;
    cpuCost : number;
    memoryCost : number;
    totalCost : number;
    existingCost : number;
    nodes : NodeDetails[];
}