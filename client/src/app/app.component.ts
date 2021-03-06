import { Component } from '@angular/core';
import { Service } from './app.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'url-shortener';
  shorten = '';

  constructor(private service: Service) {
  }

  onClick(url: string) {
    this.shorten = '';
    this.service.shorten(url)
      .subscribe(
        response => this.shorten = response.shortenUrl,
        err => this.service.showSnackbar(err.error, null, 5000),
        () => console.log('posted!')
      );
  }
}
