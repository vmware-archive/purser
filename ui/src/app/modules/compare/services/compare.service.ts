import { Injectable } from '@angular/core';

import { BACKEND_URL } from '../../../app.component'
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { CloudRegion } from '../components/compare-clouds/cloud-region';
import { CloudDetails } from '../components/compare-clouds/cloud-details';

@Injectable()

export class CompareService{

    regions : CloudRegion[] = [];
    cloudDetails : CloudDetails[] = [];
    
    getRegionsUrl = BACKEND_URL + "clouds/regions";
    postCloudRegion = BACKEND_URL + "clouds/compare";

    constructor(private http: HttpClient){ }

    getRegions() : Observable<any>{
        return this.http.get<any>(this.getRegionsUrl);     
    }

    //returns all the clouds details
    sendCloudRegion(sendCloudRegion):Observable<any>{
      return this.http.post<any>(this.postCloudRegion,sendCloudRegion);
    }
    
}