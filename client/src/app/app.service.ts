import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { tap, catchError } from 'rxjs/operators';
import { UrlData, ShortenData } from './app.model';
import { MatSnackBar } from '@angular/material';

@Injectable({
  providedIn: 'root',
})
export class Service {

  constructor(
    private http: HttpClient,
    private snackbar: MatSnackBar) { }

  postUrl = '/api/post';
  httpOptions = {
    headers: new HttpHeaders({ 'Content-Type': 'application/json' })
  };

  shorten(url: string): Observable<ShortenData> {
    const request: UrlData = { url };
    return this.http.post<ShortenData>(this.postUrl, request, this.httpOptions);
  }


  showSnackbar(message: string, action: string, duration: number) {
    this.snackbar.open(message, action, { duration });
  }

  // private handleError<T>(operation = 'operation', result?: T) {
  //   return (error: any): Observable<T> => {

  //     // TODO: send the error to remote logging infrastructure
  //     console.error(error); // log to console instead


  //     // Let the app keep running by returning an empty result.
  //     return of(result as T);
  //   };
  // }
}
